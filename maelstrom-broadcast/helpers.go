package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func readIntoMap(msg maelstrom.Message) (map[string]any, error) {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return make(map[string]any), err
	}
	return body, nil
}

func readNeighbors(nodeId string, body map[string]any) ([]string, error) {
	topology, ok := body["topology"].(map[string]any)
	if !ok {
		return nil, ErrInvalidMessage
	}
	neighbors, ok := topology[nodeId].([]any)
	if !ok {
		return nil, ErrInvalidMessage
	}
	result := make([]string, 0, len(neighbors))
	for _, value := range neighbors {
		node, ok := value.(string)
		if !ok {
			return nil, ErrInvalidMessage
		}

		result = append(result, node)
	}
	return result, nil
}

func readBroadcastMessage(body map[string]any) int {
	return int(body["message"].(float64))
}

func shouldReply(body map[string]any) bool {
	_, exists := body["msg_id"]
	return exists
}
