package main

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const LeaderKey = "cluster#leader"

func becomeLeader(storage *maelstrom.KV, nodeId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return storage.CompareAndSwap(ctx, LeaderKey, nodeId, nodeId, true)
}

func getLeader(storage *maelstrom.KV) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	value, err := storage.Read(ctx, LeaderKey)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}
