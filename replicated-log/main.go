package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	nodeLog := ReplicatedLog{}

	node.Handle("send", func(msg maelstrom.Message) error {
		req, err := parseMessage[SendRequest](msg)
		if err != nil {
			return err
		}

		if !nodeLog.Has(req.Key) {
			nodeLog.Init(req.Key)
		}

		lastOffset, err := nodeLog.Append(req.Key, req.Msg)
		if err != nil {
			return err
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
			result[key] = messageList
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
