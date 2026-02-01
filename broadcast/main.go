package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/sync/errgroup"
)

func initLogger() *log.Logger {
	file, err := os.OpenFile("/tmp/maelstrom-broadcast.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "[DEBUG] ", log.LstdFlags)
	return logger
}

func logWithCtx(l *log.Logger, nodeId string, msg string) {
	l.Printf("nodeID=%s: %s", nodeId, msg)
}

func main() {
	node := maelstrom.NewNode()
	logger := initLogger()
	state := newNodeState()

	node.Handle("topology", func(msg maelstrom.Message) error {
		parsedMsg, err := parseMessage[TopologyRequest](msg)
		if err != nil {
			return err
		}

		neighbors, ok := parsedMsg.Topology[node.ID()]
		if !ok || len(neighbors) == 0 {
			panic(fmt.Sprintf("Couldn't find node (ID:%s) neighbors.", node.ID()))
		}

		state.setNeighbors(neighbors)
		response := map[string]any{
			"type": "topology_ok",
		}

		logWithCtx(logger, node.ID(), fmt.Sprintf("Set neighbors: %v", neighbors))
		return node.Reply(msg, response)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		currentState := state.getValues()
		body := map[string]any{
			"type":     "read_ok",
			"messages": currentState,
		}
		logWithCtx(logger, node.ID(), fmt.Sprintf("Read state: %v", currentState))
		return node.Reply(msg, body)
	})

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		body, err := parseMessage[BroadcastRequest](msg)
		if err != nil {
			return err
		}

		ok := state.isExist(body.Message)
		if ok {
			logWithCtx(logger, node.ID(), fmt.Sprintln("Message already exists"))
			return node.Reply(msg, map[string]any{})
		}
		state.addValue(body.Message)
		logWithCtx(logger, node.ID(), fmt.Sprintf("Broadcast message: %v", body.Message))

		group := new(errgroup.Group)
		retry := exponentialBackoff(500*time.Millisecond, 15)
		for _, neighbor := range state.neighbors {
			// Skip sending message back to a sender
			if neighbor == msg.Src {
				continue
			}

			group.Go(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				request := func() error {
					_, err := node.SyncRPC(ctx, neighbor, body)
					return err
				}

				return retry(request, func(err error) {
					logWithCtx(logger, node.ID(), fmt.Sprintf("Error sending message to neighbor %s: %v", neighbor, err))
				})
			})
		}

		if err := group.Wait(); err != nil {
			return err
		}

		return node.Reply(msg, map[string]any{
			"msg_id": body.MsgID,
			"type":   "broadcast_ok",
		})
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
