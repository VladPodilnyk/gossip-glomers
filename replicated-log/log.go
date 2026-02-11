package main

import (
	"errors"
	"fmt"
	"sync"
)

// TODO: persist WAL on the disk
type ReplicatedLog struct {
	mu   sync.Mutex // used only during initialization
	wals map[string]*WAL
}

func (lg *ReplicatedLog) Init(key string) error {
	return nil
}

func (lg *ReplicatedLog) Has(key string) bool {
	_, ok := lg.wals[key]
	return ok
}

func (lg *ReplicatedLog) Append(key string, value int) (int, error) {
	if !lg.Has(key) {
		return -1, errors.New(fmt.Sprintf("Failed to append to the log partition %s", key))
	}
	return lg.wals[key].Append(value)
}

func (lg *ReplicatedLog) Read(key string, offset uint, limit uint) ([][]int, error) {
	if limit == 0 || offset > uint(lg.wals[key].Offset()) {
		return [][]int{}, nil
	}
	return lg.wals[key].Read(offset, limit)
}

func (lg *ReplicatedLog) Commit(offsets map[string]int) error {
	for key, offset := range offsets {
		// Ignore incorrect offsets - should be fine for the toy implementation
		if offset > lg.wals[key].LastCommittedOffset() {
			err := lg.wals[key].Commit(offset)
			if err != nil {
				return err
			}
		}
	}
}

func (lg *ReplicatedLog) GetCommittedOffsets(keys []string) map[string]int {
	result := make(map[string]int, len(keys))
	for _, key := range keys {
		result[key] = lg.wals[key].LastCommittedOffset()
	}
	return result
}
