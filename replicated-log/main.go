package main

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	var leaderNode string
	nodeLog := newReplicatedLog()
	kv := maelstrom.NewLinKV(node)

	node.Handle("init", func(msg maelstrom.Message) error {
		// register a node as a leader or fetch current leader
		leaderNode = node.ID()
		err := becomeLeader(kv, leaderNode)
		if err != nil {
			leaderNode, err = getLeader(kv)
			if err != nil {
				panic("Failed to assign leader")
			}
		}
		return nil
	})

	node.Handle("send", func(msg maelstrom.Message) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := parseMessage[SendRequest](msg)
		if err != nil {
			return err
		}

		// if not a leader then proxy request to a leader node
		if msg.Src != leaderNode && leaderNode != node.ID() {
			rsp, err := node.SyncRPC(ctx, leaderNode, msg.Body)
			if err != nil {
				return err
			}
			return node.Reply(msg, rsp.Body)
		}

		if err := nodeLog.InitIfNotExists(node.ID(), req.Key); err != nil {
			return err
		}

		lastOffset, err := nodeLog.Append(req.Key, req.Msg)
		if err != nil {
			return err
		}

		if leaderNode != node.ID() {
			// super simple commit: wait for all nodes to write a message.
			for _, id := range node.NodeIDs() {
				if node.ID() != id {
					rsp, err := node.SyncRPC(ctx, id, msg.Body)
					if err != nil {
						return err
					}
					req, err := parseMessage[SendRequest](rsp)
					if err != nil || req.Type != "send_ok" {
						return err
					}
				}
			}
		}

		return node.Reply(msg, map[string]any{
			"msg_id": req.MsgId,
			"type":   "send_ok",
			"offset": lastOffset,
		})
	})

	node.Handle("poll", func(msg maelstrom.Message) error {
		req, err := parseMessage[PollRequest](msg)
		if err != nil {
			return err
		}

		result := map[string][][]int{}
		for key, offset := range req.Offsets {
			messageList, err := nodeLog.Read(key, uint(offset), 5)
			if err != nil {
				return err
			}

			if len(messageList) > 0 {
				result[key] = messageList
			}
		}

		return node.Reply(msg, map[string]any{
			"msg_id": req.MsgId,
			"type":   "poll_ok",
			"msgs":   result,
		})
	})

	node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		req, err := parseMessage[CommitRequest](msg)
		if err != nil {
			return err
		}

		nodeLog.Commit(req.Offsets)
		return node.Reply(msg, map[string]any{
			"msg_id": req.MsgId,
			"type":   "commit_offsets_ok",
		})
	})

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		req, err := parseMessage[ListCommitedOffsetsRequest](msg)
		if err != nil {
			return err
		}

		result := nodeLog.GetCommittedOffsets(req.Keys)

		return node.Reply(msg, map[string]any{
			"msg_id":  req.MsgId,
			"type":    "list_committed_offsets_ok",
			"offsets": result,
		})
	})

	if err := node.Run(); err != nil {
		panic(err)
	}
}
