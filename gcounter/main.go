package main

import (
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	fmt.Printf("Node ID: %s\n", node.ID())
}
