package main

import (
	"sync"
)

type ReplicatedLog struct {
	mu   sync.RWMutex
	wals map[string]*WAL
}

func newReplicatedLog() *ReplicatedLog {
	return &ReplicatedLog{
		wals: make(map[string]*WAL),
	}
}

func (lg *ReplicatedLog) InitIfNotExists(nodeId string, key string) error {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	if _, ok := lg.wals[key]; ok {
		return nil
	}

	wal, err := newWAL(nodeId, key)
	if err != nil {
		return err
	}
	lg.wals[key] = wal
	return nil
}

func (lg *ReplicatedLog) Has(key string) bool {
	lg.mu.RLock()
	defer lg.mu.RUnlock()

	_, ok := lg.wals[key]
	return ok
}

func (lg *ReplicatedLog) Append(key string, value int) (int, error) {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	return lg.wals[key].Append(value)
}

func (lg *ReplicatedLog) Read(key string, offset uint, limit uint) ([][]int, error) {
	lg.mu.RLock()
	defer lg.mu.RUnlock()

	if _, ok := lg.wals[key]; !ok {
		return [][]int{}, nil
	}

	if limit == 0 || offset > uint(lg.wals[key].Offset()) {
		return [][]int{}, nil
	}
	return lg.wals[key].Read(offset, limit)
}

func (lg *ReplicatedLog) Commit(offsets map[string]int) error {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	for key, offset := range offsets {
		err := lg.wals[key].Commit(offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (lg *ReplicatedLog) GetCommittedOffsets(keys []string) map[string]int {
	lg.mu.RLock()
	defer lg.mu.RUnlock()

	result := make(map[string]int, len(keys))
	for _, key := range keys {
		if _, ok := lg.wals[key]; ok && lg.wals[key].LastCommittedOffset() >= 0 {
			result[key] = lg.wals[key].LastCommittedOffset()
		}
	}
	return result
}
