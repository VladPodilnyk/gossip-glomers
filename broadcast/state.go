package main

import (
	"slices"
	"sync"
)

type NodeState struct {
	mu sync.Mutex
	// TODO: use Sets
	values    []float64
	neighbors []string
}

func (ns *NodeState) addValue(value float64) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.values = append(ns.values, value)
}

func (ns *NodeState) getValues() []float64 {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	values := make([]float64, len(ns.values))
	copy(values, ns.values)

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
	return slices.Contains(ns.values, value)
}
