package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type LogEntry struct {
	Offset int `json:"offset"`
	Data   int `json:"data"`
}

type Metadata struct {
	Offset              int `json:"offset"`
	LastCommittedOffset int `json:"last_committed_offset"`
}

type WAL struct {
	mu                  sync.RWMutex
	path                string
	metadata_path       string
	file                *os.File
	metadata            *os.File
	offset              int
	lastCommittedOffset int
}

func newWAL(nodeId string, key string) (*WAL, error) {
	path := fmt.Sprintf("/tmp/maelstrom/log/%s/%s/wal.log", nodeId, key)
	metadata_path := fmt.Sprintf("/tmp/maelstrom/log/%s/%s/wal.metadata", nodeId, key)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}

	metadata, err := os.OpenFile(metadata_path, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}

	return &WAL{
		path:                path,
		metadata_path:       metadata_path,
		file:                file,
		metadata:            metadata,
		offset:              -1,
		lastCommittedOffset: -1,
	}, nil
}

func (w *WAL) Append(data int) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	newOffset := w.offset + 1
	entry := LogEntry{
		Offset: newOffset,
		Data:   data,
	}
	updatedMetadata := Metadata{
		Offset:              newOffset,
		LastCommittedOffset: w.lastCommittedOffset,
	}

	logEntryRaw, err := json.Marshal(entry)
	if err != nil {
		return -1, err
	}

	metadataRaw, err := json.Marshal(updatedMetadata)
	if err != nil {
		return -1, err
	}

	entryBytes := append(logEntryRaw, '\n')
	metadataBytes := append(metadataRaw, '\n')

	if _, err := w.file.Write(entryBytes); err != nil {
		return -1, err
	}

	if _, err := w.metadata.Seek(0, io.SeekStart); err != nil {
		return -1, err
	}

	if _, err := w.metadata.Write(metadataBytes); err != nil {
		return -1, err
	}

	w.offset = newOffset
	return w.offset - 1, nil
}

func (w *WAL) Commit(offset int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if offset < w.lastCommittedOffset || offset > w.offset {
		return errors.New("invalid offset")
	}

	w.lastCommittedOffset = offset
	updatedMetadata := Metadata{
		Offset:              w.offset,
		LastCommittedOffset: w.lastCommittedOffset,
	}

	metadataRaw, err := json.Marshal(updatedMetadata)
	if err != nil {
		return err
	}

	if _, err := w.metadata.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := w.metadata.Write(metadataRaw); err != nil {
		return err
	}

	return nil
}

func (w *WAL) Read(offset uint, limit uint) ([][]int, error) {
	scanner := bufio.NewScanner(w.file)
	currOffset := 0

	result := make([][]int, 0, limit)
	for scanner.Scan() {
		if currOffset >= int(offset) {
			entry := LogEntry{}
			err := json.Unmarshal(scanner.Bytes(), &entry)
			if err != nil {
				return nil, err
			}
			result = append(result, []int{entry.Offset, entry.Data})
			if len(result) == int(limit) {
				break
			}
		}
		currOffset += 1
	}
	return result, nil
}

func (w *WAL) Offset() int {
	return w.offset
}

func (w *WAL) LastCommittedOffset() int {
	return w.lastCommittedOffset
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.file.Close(); err != nil {
		return err
	}

	if err := w.metadata.Close(); err != nil {
		return err
	}

	return nil
}
