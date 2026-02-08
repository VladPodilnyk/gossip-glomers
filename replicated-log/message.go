package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type SendRequest struct {
	MsgId int    `json:"msg_id"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Msg   int    `json:"msg"`
}

type PollRequest struct {
	MsgId   int            `json:"msg_id"`
	Type    string         `json:"type"`
	Key     string         `json:"key"`
	Offsets map[string]int `json:"offsets"`
}

type CommitRequest struct {
	MsgId   int            `json:"msg_id"`
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

type ListCommitedOffsetsRequest struct {
	MsgId int      `json:"msg_id"`
	Type  string   `json:"type"`
	Keys  []string `json:"keys"`
}

func parseMessage[T any](data maelstrom.Message) (T, error) {
	var msg T
	if err := json.Unmarshal(data.Body, &msg); err != nil {
		return msg, err
	}
	return msg, nil
}
