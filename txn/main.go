package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func broadcast(node *maelstrom.Node, msg map[string]any) {
	for _, id := range node.NodeIDs() {
		if id == node.ID() {
			continue
		}

		if err := rpcWithRetry(node, id, msg, 5); err != nil {
			// Fail fast
			panic(err)
		}
	}
}

func main() {
	node := maelstrom.NewNode()
	state := newState()

	node.Handle("txn", func(msg maelstrom.Message) error {
		txnMessage, err := parseMessage[TxnMessage](msg)
		if err != nil {
			return err
		}

		results := make([][]any, 0, len(txnMessage.Txn))
		state.mu.Lock()
		defer state.mu.Unlock()

		writes := make(map[int]int)
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
				writes[key] = *value
				results = append(results, []any{opType, key, value})
			}
		}

		// Broadcast (simply fire&forget - works for this toy-like solution)
		go broadcast(node, map[string]any{
			"type":   "sync",
			"writes": writes,
		})

		return node.Reply(msg, map[string]any{
			"type":   "txn_ok",
			"msg_id": txnMessage.MsgId,
			"txn":    results,
		})
	})

	node.Handle("sync", func(msg maelstrom.Message) error {
		syncMessage, err := parseMessage[SyncMessage](msg)
		if err != nil {
			return err
		}

		state.mu.Lock()
		defer state.mu.Unlock()

		for k, v := range syncMessage.Writes {
			state.store[k] = v
		}

		return node.Reply(msg, map[string]any{
			"type": "sync_ok",
		})
	})

	if err := node.Run(); err != nil {
		panic(err)
	}
}
