package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()

	node.Handle("send", func(msg maelstrom.Message) error {
		return nil
	})

	node.Handle("poll", func(msg maelstrom.Message) error {
		return nil
	})

	node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		return nil
	})

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		return nil
	})

	if err := node.Run(); err != nil {
		panic(err)
	}
}
