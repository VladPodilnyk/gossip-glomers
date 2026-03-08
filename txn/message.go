package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TxnMessage struct {
	MsgId int     `json:"msg_id"`
	Type  string  `json:"type"`
	Txn   [][]any `json:"txn"`
}

func parseMessage(data maelstrom.Message) (TxnMessage, error) {
	var msg TxnMessage
	if err := json.Unmarshal(data.Body, &msg); err != nil {
		return msg, err
	}
	return msg, nil
}

func extractOpDetails(op []any) (string, int, *int) {
	opType := op[0].(string)
	key := op[1].(float64)

	var value *int = nil
	if op[2] != nil {
		fValue := op[2].(float64)
		iValue := int(fValue)
		value = &iValue
	}
	return opType, int(key), value
}
