package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	state := newState()

	node.Handle("txn", func(msg maelstrom.Message) error {
		txnMessage, err := parseMessage(msg)
		if err != nil {
			return err
		}

		results := make([][]any, 0, len(txnMessage.Txn))
		state.mu.Lock()
		defer state.mu.Unlock()

		for _, op := range txnMessage.Txn {
			opType, key, value := extractOpDetails(op)
			if opType == "r" {
				storedValue, ok := state.Read(key)
				if ok {
					results = append(results, []any{opType, key, storedValue})
				} else {
					results = append(results, []any{opType, key, nil})
				}
			} else {
				state.Write(key, *value)
				results = append(results, []any{opType, key, value})
			}
		}

		return node.Reply(msg, map[string]any{
			"type":   "txn_ok",
			"msg_id": txnMessage.MsgId,
			"txn":    results,
		})
	})

	if err := node.Run(); err != nil {
		panic(err)
	}
}
