package main

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func rpcRequest(node *maelstrom.Node, dest string, payload any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := node.SyncRPC(ctx, dest, payload)
	return err
}

func rpcWithRetry(node *maelstrom.Node, dest string, payload any, retries uint) error {
	var err error
	for i := 0; i < int(retries); i++ {
		if err = rpcRequest(node, dest, payload); err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return err
}
