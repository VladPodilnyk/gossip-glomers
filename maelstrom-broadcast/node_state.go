package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

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

func (s *nodeState) GetHistory() []int {
	keys := make([]int, 0, len(s.seen))
	for k := range s.seen {
		keys = append(keys, k)
	}
	return keys
}

func (s *nodeState) Broadcast(message int) {
	body := make(map[string]any)
	body["message"] = message

	for _, n := range s.neighbors {
		// Fire and forget
		go func() {
			s.node.Send(n, body)
		}()
	}
}
