package main

import "sync"

type State struct {
	mu    sync.Mutex
	store map[int]int
}

func newState() *State {
	return &State{store: make(map[int]int)}
}

func (s *State) Read(key int) (int, bool) {
	v, ok := s.store[key]
	return v, ok
}

func (s *State) Write(key int, value int) {
	s.store[key] = value
}
