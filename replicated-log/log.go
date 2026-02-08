package main

import "sync"

type LogEntry struct {
	Offset int
	Data   int
}

// TODO: persist WAL on the disk
type ReplicatedLog struct {
	mu                  sync.Mutex
	LastOffset          map[string]int
	LastCommittedOffset map[string]int
	Messages            map[string][]LogEntry // WAL
}

func newReplicatedLog() *ReplicatedLog {
	return &ReplicatedLog{
		LastOffset:          make(map[string]int),
		LastCommittedOffset: make(map[string]int),
		Messages:            make(map[string][]LogEntry),
	}
}

// TODO: protect my lock for multi-node setup
func (lg *ReplicatedLog) Append(key string, value int) {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	lastOffset, ok := lg.LastOffset[key]
	if !ok {
		lastOffset = -1
	}
	newOffset := lastOffset + 1
	entry := LogEntry{
		Offset: newOffset,
		Data:   value,
	}
	lg.Messages[key] = append(lg.Messages[key], entry)
	lg.LastOffset[key] = newOffset
}

func (lg *ReplicatedLog) ReadMessages(key string, offset uint, limit uint) [][]int {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	if limit == 0 || offset > uint(len(lg.Messages[key])) {
		return [][]int{}
	}
	result := make([][]int, 0, limit)
	for i := offset; i < offset+limit; i++ {
		if i >= uint(len(lg.Messages[key])) {
			break
		}
		result = append(result, []int{lg.Messages[key][i].Offset, lg.Messages[key][i].Data})
	}
	return result
}

func (lg *ReplicatedLog) Commit(offsets map[string]int) {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	for key, offset := range offsets {
		if offset > lg.LastCommittedOffset[key] {
			lg.LastCommittedOffset[key] = offset
		}
	}
}

func (lg *ReplicatedLog) GetCommittedOffsets(keys []string) map[string]int {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	result := make(map[string]int, len(keys))
	for _, key := range keys {
		result[key] = lg.LastCommittedOffset[key]
	}
	return result
}
