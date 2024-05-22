package order_queue

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

func (d *Disruptor) SaveSnapshotToFile() error {
	snapshotDir := os.Getenv("SNAPSHOT_DIR")
	if snapshotDir == "" {
		return errors.New("SNAPSHOT_DIR environment variable is not set")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(d.buffer)
	if err != nil {
		return err
	}
	err = encoder.Encode(d.writeCursor)
	if err != nil {
		return err
	}
	err = encoder.Encode(d.readCursor)
	if err != nil {
		return err
	}

	snapshotPath := filepath.Join(snapshotDir, fmt.Sprintf("snapshot-%d.gob", atomic.LoadInt64(&d.writeCursor)))
	err = os.WriteFile(snapshotPath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (d *Disruptor) LoadSnapshotFromFile() error {
	snapshotDir := os.Getenv("SNAPSHOT_DIR")
	if snapshotDir == "" {
		return errors.New("SNAPSHOT_DIR environment variable is not set")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	entries, err := os.ReadDir(snapshotDir)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return errors.New("no snapshots found in the snapshot directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			snapshotPath := filepath.Join(snapshotDir, entry.Name())
			file, err := os.Open(snapshotPath)
			if err != nil {
				continue
			}
			defer file.Close()

			decoder := gob.NewDecoder(file)
			err = decoder.Decode(&d.buffer)
			if err != nil {
				return err
			}
			err = decoder.Decode(&d.writeCursor)
			if err != nil {
				return err
			}
			err = decoder.Decode(&d.readCursor)
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("no available snapshots found in the snapshot directory")
}
