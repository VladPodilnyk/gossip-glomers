package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TopologyRequest struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
	MsgID    int                 `json:"msg_id"`
}

type BroadcastRequest struct {
	Type    string  `json:"type"`
	MsgID   int     `json:"msg_id"`
	Message float64 `json:"message"`
}

func parseMessage[T any](data maelstrom.Message) (T, error) {
	var msg T
	if err := json.Unmarshal(data.Body, &msg); err != nil {
		return msg, err
	}
	return msg, nil
}
