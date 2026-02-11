package main

import (
	"os"
	"sync"
)

type LogEntry struct {
	offset int
	data   int
}

type WAL struct {
	mu                  sync.Mutex
	path                string
	file                *os.File
	offset              int
	lastCommittedOffset int
}

func newWAL() *WAL {
	return &WAL{offset: -1, lastCommittedOffset: -1}
}

func (w *WAL) Append(data int) (int, error) {
	return 0, nil
}

func (w *WAL) Commit(offset int) error {
	return nil
}

func (w *WAL) Read(offset uint, limit uint) ([][]int, error) {
	return nil, nil
}

func (w *WAL) Offset() int {
	return w.offset
}

func (w *WAL) LastCommittedOffset() int {
	return w.lastCommittedOffset
}

func (w *WAL) Close() error {
	return nil
}

func (w *WAL) recovery() error {
	return nil
}
