package main

import (
	"fmt"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// TODO refactor code (list)
// 1. Improve parsing of messages
// 2. Implement 3c handling network partitions
// 3. Decrease code duplication and move common logic to separate functions and files.

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
	state := &NodeState{}

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
		// TODO: rewrite with goroutines and exponential backoff to handle failures
		for _, neighbor := range state.neighbors {
			// Skip sending message back to a sender
			if neighbor == msg.Src {
				continue
			}

			err := node.Send(neighbor, body)
			if err != nil {
				logWithCtx(logger, node.ID(), fmt.Sprintf("Error sending message to neighbor %s: %v", neighbor, err))
			}
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
