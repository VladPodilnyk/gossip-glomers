package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// TODO refactor code (list)
// 1. Improve parsing of messages
// 2. Implement 3c handling network partitions
// 3. Decrease code duplication and move common logic to separate functions and files.

type topologyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
	MsgID    int                 `json:"msg_id"`
}

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

type nodeState struct {
	mu sync.Mutex
	// TODO: use Sets
	values    []float64
	neighbors []string
}

func (ns *nodeState) addValue(value float64) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.values = append(ns.values, value)
}

func (ns *nodeState) getValues() []float64 {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	values := make([]float64, len(ns.values))
	copy(values, ns.values)

	return values
}

// Since we get "topology" request only once during cluster initialization
// we set values without locking
func (ns *nodeState) setNeighbors(neighbors []string) {
	ns.neighbors = neighbors
}

func (ns *nodeState) isExist(value float64) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	return slices.Contains(ns.values, value)
}

func main() {
	node := maelstrom.NewNode()
	logger := initLogger()
	state := &nodeState{}

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body topologyMsg
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		neighbors, ok := body.Topology[node.ID()]
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
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		currentState := state.getValues()
		body["messages"] = currentState

		logWithCtx(logger, node.ID(), fmt.Sprintf("Read state: %v", currentState))
		return node.Reply(msg, body)
	})

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		value := body["message"].(float64)
		ok := state.isExist(value)
		if ok {
			logWithCtx(logger, node.ID(), fmt.Sprintln("Message already exists"))
			delete(body, "message")
			body["type"] = "broadcast_ok"
			return node.Reply(msg, body)
		}
		state.addValue(value)

		logWithCtx(logger, node.ID(), fmt.Sprintf("Broadcast message: %v", value))
		// TODO: rewrite with goroutines and exponential backoff to handle failures
		for _, neighbor := range state.neighbors {
			// Skip sending message back to a sender
			if neighbor == msg.Src {
				continue
			}

			body["type"] = "broadcast"
			body["message"] = value
			err := node.Send(neighbor, body)
			if err != nil {
				logWithCtx(logger, node.ID(), fmt.Sprintf("Error sending message to neighbor %s: %v", neighbor, err))
			}
		}
		delete(body, "message")
		body["type"] = "broadcast_ok"
		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
