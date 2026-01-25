package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var counter atomic.Int64
	node := maelstrom.NewNode()

	node.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		generatedId := fmt.Sprintf("%s:%d", node.ID(), counter.Add(1))
		body["type"] = "generate_ok"
		body["id"] = generatedId
		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
