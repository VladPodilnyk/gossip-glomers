package main

import (
	"strconv"
	"sync"
)

type NodeState struct {
	mu        sync.Mutex
	values    map[string]bool
	neighbors []string
}

func newNodeState() *NodeState {
	return &NodeState{
		values:    make(map[string]bool),
		neighbors: make([]string, 0),
	}
}

func (ns *NodeState) addValue(value float64) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	key := strconv.FormatFloat(value, 'f', -1, 64)
	ns.values[key] = true
}

func (ns *NodeState) getValues() []float64 {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	values := make([]float64, len(ns.values))
	i := 0
	for key := range ns.values {
		values[i], _ = strconv.ParseFloat(key, 64)
		i++
	}

	return values
}

// Since we get "topology" request only once during cluster initialization
// we set values without locking
func (ns *NodeState) setNeighbors(neighbors []string) {
	ns.neighbors = neighbors
}

func (ns *NodeState) isExist(value float64) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	key := strconv.FormatFloat(value, 'f', -1, 64)
	return ns.values[key]
}
