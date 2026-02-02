package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func initLogger() *log.Logger {
	file, err := os.OpenFile("/tmp/maelstrom-broadcast.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "[DEBUG] ", log.LstdFlags)
	return logger
}

func logWithCtx(l *log.Logger, nodeId string, msg string) {
	l.Printf("nodeID=%s: %s", nodeId, msg)
}

type AddMessage struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
	MsgID int    `json:"msg_id"`
}

func main() {
	node := maelstrom.NewNode()
	logger := initLogger()
	kv := maelstrom.NewSeqKV(node)

	node.Handle("add", func(msg maelstrom.Message) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var addMessage AddMessage
		if err := json.Unmarshal(msg.Body, &addMessage); err != nil {
			return err
		}

		currValue, err := kv.ReadInt(ctx, node.ID())
		if err != nil && maelstrom.ErrorCode(err) != maelstrom.KeyDoesNotExist {
			logWithCtx(logger, node.ID(), fmt.Sprintf("error reading key: %v, nodeID=%s", err, node.ID()))
			return err
		}

		if maelstrom.ErrorCode(err) == maelstrom.KeyDoesNotExist {
			currValue = 0
		}

		kv.CompareAndSwap(ctx, node.ID(), currValue, currValue+addMessage.Delta, true)
		return node.Reply(msg, map[string]any{
			"msg_id": addMessage.MsgID,
			"type":   "add_ok",
		})
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		gcounterValue := 0
		for _, key := range node.NodeIDs() {
			value, err := kv.ReadInt(ctx, key)
			if err != nil {
				if maelstrom.ErrorCode(err) == maelstrom.KeyDoesNotExist {
					continue
				}
				logWithCtx(logger, node.ID(), fmt.Sprintf("error reading key: %v, nodeID=%s", err, node.ID()))
				return err
			}
			gcounterValue += value
		}

		body["type"] = "read_ok"
		body["value"] = gcounterValue
		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		panic(err)
	}
}
