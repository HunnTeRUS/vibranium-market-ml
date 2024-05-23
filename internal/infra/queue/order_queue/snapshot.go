package order_queue

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"os"
	"path/filepath"
	"sync/atomic"
)

func (d *Disruptor) SaveSnapshotToFile() error {
	if len(d.queue) == 0 {
		return nil
	}

	snapshotDir := os.Getenv("SNAPSHOT_DIR")
	if snapshotDir == "" {
		msg := "SNAPSHOT_DIR environment variable is not set"

		logger.Warn(msg)
		return errors.New(msg)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	var orders []*order.Order
	for len(d.queue) > 0 {
		order := <-d.queue
		orders = append(orders, order)
	}
	err := encoder.Encode(orders)
	if err != nil {
		logger.Warn(err.Error())
		return err
	}
	err = encoder.Encode(atomic.LoadInt64(&d.writeCursor))
	if err != nil {
		logger.Warn(err.Error())
		return err
	}
	err = encoder.Encode(atomic.LoadInt64(&d.readCursor))
	if err != nil {
		logger.Warn(err.Error())
		return err
	}

	snapshotPath := filepath.Join(snapshotDir, fmt.Sprintf("snapshot-%d.gob", atomic.LoadInt64(&d.writeCursor)))
	err = os.WriteFile(snapshotPath, buffer.Bytes(), 0644)
	if err != nil {
		logger.Warn(err.Error())
		return err
	}

	for _, order := range orders {
		d.queue <- order
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
		logger.Warn(err.Error())
		return err
	}

	if len(entries) == 0 {
		msg := "no snapshots found in the snapshot directory"
		logger.Warn(msg)
		return errors.New(msg)
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
			var orders []*order.Order
			err = decoder.Decode(&orders)
			if err != nil {
				return err
			}
			var writeCursor, readCursor int64
			err = decoder.Decode(&writeCursor)
			if err != nil {
				return err
			}
			err = decoder.Decode(&readCursor)
			if err != nil {
				return err
			}

			d.queue = make(chan *order.Order, cap(d.queue))
			for _, order := range orders {
				d.queue <- order
			}
			atomic.StoreInt64(&d.writeCursor, writeCursor)
			atomic.StoreInt64(&d.readCursor, readCursor)

			return nil
		}
	}

	msg := "no available snapshots found in the snapshot directory"
	logger.Warn(msg)
	return errors.New(msg)
}
