package main

import (
	"fmt"
	"sync"
)

// TODO: persist WAL on the disk
type ReplicatedLog struct {
	mu   sync.RWMutex
	wals map[string]*WAL
}

func newReplicatedLog() *ReplicatedLog {
	return &ReplicatedLog{
		wals: make(map[string]*WAL),
	}
}

func (lg *ReplicatedLog) Init(nodeId string, key string) error {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	wal, err := newWAL(nodeId, key)
	if err != nil {
		return err
	}
	lg.wals[key] = wal
	return nil
}

func (lg *ReplicatedLog) Has(key string) bool {
	_, ok := lg.wals[key]
	return ok
}

func (lg *ReplicatedLog) Append(key string, value int) (int, error) {
	if !lg.Has(key) {
		return -1, fmt.Errorf("Failed to append to the log partition %s", key)
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
		err := lg.wals[key].Commit(offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (lg *ReplicatedLog) GetCommittedOffsets(keys []string) map[string]int {
	result := make(map[string]int, len(keys))
	for _, key := range keys {
		if lg.Has(key) {
			result[key] = lg.wals[key].LastCommittedOffset()
		}
	}
	return result
}
