package main

import (
	"errors"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var ErrInvalidMessage = errors.New("invalid message")

func handleBroadcast(state *nodeState) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		body, err := readIntoMap(msg)
		if err != nil {
			return err
		}

		data := readBroadcastMessage(body)
		_, shouldAbortBroadcasting := state.seen[data]
		state.AddMessage(data)

		if !shouldAbortBroadcasting {
			state.Broadcast(data)
		}

		response := make(map[string]any)
		response["type"] = "broadcast_ok"
		return state.node.Reply(msg, response)
	}
}

func handleRead(state *nodeState) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		body, err := readIntoMap(msg)
		if err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["messages"] = state.GetHistory()
		return state.node.Reply(msg, body)
	}
}

func handleTopology(state *nodeState) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		body, err := readIntoMap(msg)
		if err != nil {
			return err
		}

		state.neighbors, err = readNeighbors(state.node.ID(), body)
		if err != nil {
			return err
		}

		response := make(map[string]any)
		response["type"] = "topology_ok"
		return state.node.Reply(msg, response)
	}
}

func main() {
	node := maelstrom.NewNode()
	state := newNodeState(node)

	node.Handle("broadcast", handleBroadcast(state))
	node.Handle("read", handleRead(state))
	node.Handle("topology", handleTopology(state))

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
