package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func initLogger() *log.Logger {
	file, err := os.OpenFile("/tmp/broadcast.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	var state []float64
	var neighbors []string

	// Not required for the part 3a
	node.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology := body["topology"].(map[string]any)
		neighbors, _ = topology[node.ID()].([]string)
		if neighbors == nil {
			neighbors = []string{}
		}

		response := map[string]any{
			"type": "topology_ok",
		}

		logWithCtx(logger, node.ID(), fmt.Sprintf("Set neighbors: %v", neighbors))
		return node.Reply(msg, response)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		// TODO: Works for the part 3a, but won't work for multiple nodes
		body["messages"] = state

		logWithCtx(logger, node.ID(), fmt.Sprintf("Read state: %v", state))
		return node.Reply(msg, body)
	})

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// TODO: Will not work for concurrent operations
		value := body["message"].(float64)
		state = append(state, value)

		logWithCtx(logger, node.ID(), fmt.Sprintf("Broadcast message: %v", value))
		// TODO: Not required for the part 3a
		// for _, neighbor := range neighbors {
		// 	// TODO: check if ID is attached or not
		// 	err := node.Send(neighbor, body)
		// 	if err != nil {
		// 		logWithCtx(logger, node.ID(), fmt.Sprintf("Error sending message to neighbor %s: %v", neighbor, err))
		// 	}
		// }
		delete(body, "message")
		body["type"] = "broadcast_ok"

		node.Reply(msg, body)
		return nil
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
