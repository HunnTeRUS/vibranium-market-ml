package order_queue

import (
	"bytes"
	"encoding/gob"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveSnapshot(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snapshot")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("SNAPSHOT_DIR", tmpDir)
	defer os.Unsetenv("SNAPSHOT_DIR")

	d := &Disruptor{
		queue: make(chan *order.Order, 10),
	}
	order1 := &order.Order{ID: "order1", UserID: "user1", Price: 50.0, Amount: 100}
	order2 := &order.Order{ID: "order2", UserID: "user2", Price: 100.0, Amount: 200}
	d.queue <- order1
	d.queue <- order2

	t.Run("save snapshot successfully", func(t *testing.T) {
		err := d.SaveSnapshotToFile()
		assert.NoError(t, err)

		files, err := os.ReadDir(tmpDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, files)

		for _, file := range files {
			filePath := filepath.Join(tmpDir, file.Name())
			data, err := os.ReadFile(filePath)
			assert.NoError(t, err)

			var orders []*order.Order
			var writeCursor, readCursor int64
			buffer := bytes.NewBuffer(data)
			decoder := gob.NewDecoder(buffer)
			err = decoder.Decode(&orders)
			assert.NoError(t, err)
			err = decoder.Decode(&writeCursor)
			assert.NoError(t, err)
			err = decoder.Decode(&readCursor)
			assert.NoError(t, err)

			assert.Len(t, orders, 2)
			assert.Equal(t, "order1", orders[0].ID)
			assert.Equal(t, "order2", orders[1].ID)
		}
	})

	t.Run("save snapshot with empty queue", func(t *testing.T) {
		dEmpty := &Disruptor{
			queue: make(chan *order.Order, 10),
		}
		err := dEmpty.SaveSnapshotToFile()
		assert.NoError(t, err)
	})
}

func TestLoadSnapshotFromFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snapshot")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("SNAPSHOT_DIR", tmpDir)
	defer os.Unsetenv("SNAPSHOT_DIR")

	d := &Disruptor{
		queue: make(chan *order.Order, 10),
	}
	order1 := &order.Order{ID: "order1", UserID: "user1", Price: 50.0, Amount: 100}
	order2 := &order.Order{ID: "order2", UserID: "user2", Price: 100.0, Amount: 200}
	d.queue <- order1
	d.queue <- order2
	err = d.SaveSnapshotToFile()
	assert.NoError(t, err)

	t.Run("load snapshot successfully", func(t *testing.T) {
		dLoaded := &Disruptor{
			queue: make(chan *order.Order, 10),
		}
		err := dLoaded.LoadSnapshotFromFile()
		assert.NoError(t, err)

		assert.Equal(t, cap(d.queue), cap(dLoaded.queue))

		var orders []*order.Order
		for len(dLoaded.queue) > 0 {
			orders = append(orders, <-dLoaded.queue)
		}
		assert.Len(t, orders, 2)
		assert.Equal(t, order1, orders[0])
		assert.Equal(t, order2, orders[1])
	})

	t.Run("load snapshot with no snapshots", func(t *testing.T) {
		os.RemoveAll(tmpDir)
		err := d.LoadSnapshotFromFile()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot find the file specified")
	})
}

func TestLoadSnapshotFromFile_FileNotFound(t *testing.T) {
	os.Setenv("SNAPSHOT_DIR", "nonexistent_dir")
	defer os.Unsetenv("SNAPSHOT_DIR")

	d := &Disruptor{
		queue: make(chan *order.Order, 10),
	}

	err := d.LoadSnapshotFromFile()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot find the file specified")
}

func TestSaveSnapshot_DirectoryCreationError(t *testing.T) {
	os.Setenv("SNAPSHOT_DIR", string([]rune{0}))
	defer os.Unsetenv("SNAPSHOT_DIR")

	d := &Disruptor{
		queue: make(chan *order.Order, 10),
	}
	order1 := &order.Order{ID: "order1", UserID: "user1", Price: 50.0, Amount: 100}
	d.queue <- order1

	err := d.SaveSnapshotToFile()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SNAPSHOT_DIR environment variable is not set")
}
