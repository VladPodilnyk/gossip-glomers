package main

import (
	"encoding/json"
	"errors"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var ErrInvalidMessage = errors.New("invalid message")

type nodeState struct {
	seen      map[int]bool
	neighbors []string
	node      *maelstrom.Node
}

func newNodeState(node *maelstrom.Node) *nodeState {
	return &nodeState{seen: make(map[int]bool), node: node}
}

func (s *nodeState) AddMessage(newMessage int) {
	s.seen[newMessage] = true
}

func (s *nodeState) Read() []int {
	keys := make([]int, 0, len(s.seen))
	for k := range s.seen {
		keys = append(keys, k)
	}
	return keys
}

func readIntoMap(msg maelstrom.Message) (map[string]any, error) {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return make(map[string]any), err
	}
	return body, nil
}

func readTopology(nodeId string, body map[string]any) ([]string, error) {
	topology, ok := body["topology"].(map[string]any)
	if !ok {
		return nil, ErrInvalidMessage
	}
	neighbors, ok := topology[nodeId].([]interface{})
	if !ok {
		return nil, ErrInvalidMessage
	}
	result := make([]string, 0, len(neighbors))
	for _, value := range neighbors {
		node, ok := value.(string)
		if !ok {
			return nil, ErrInvalidMessage
		}

		value = append(result, node)
	}
	return result, nil
}

func handleBroadcast(state *nodeState) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		body, err := readIntoMap(msg)
		if err != nil {
			return err
		}

		state.AddMessage(int(body["message"].(float64)))

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
		body["messages"] = state.Read()
		return state.node.Reply(msg, body)
	}
}

func handleTopology(state *nodeState) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		body, err := readIntoMap(msg)
		if err != nil {
			return err
		}

		state.neighbors, err = readTopology(state.node.ID(), body)
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
