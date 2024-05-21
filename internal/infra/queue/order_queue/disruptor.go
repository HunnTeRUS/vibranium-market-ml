package order_queue

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"sync"
	"sync/atomic"
)

var ErrBufferFull = errors.New("buffer is full")

type Disruptor struct {
	buffer      []interface{}
	writeCursor int64
	readCursor  int64
	mu          sync.Mutex
}

func NewDisruptor(initialSize int) *Disruptor {
	d := &Disruptor{
		buffer: make([]interface{}, initialSize),
	}

	err := d.LoadSnapshotFromS3()
	if err != nil {
		logger.Warn("No previous snapshot found or failed to load")
	}

	return d
}

func (d *Disruptor) Enqueue(value interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	nextWriteCursor := (atomic.LoadInt64(&d.writeCursor) + 1) % int64(len(d.buffer))
	if nextWriteCursor == atomic.LoadInt64(&d.readCursor) {
		d.expandBuffer()
	}

	seq := atomic.AddInt64(&d.writeCursor, 1) % int64(len(d.buffer))
	d.buffer[seq] = value
	return nil
}

func (d *Disruptor) Dequeue() interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()

	if atomic.LoadInt64(&d.readCursor) == atomic.LoadInt64(&d.writeCursor) {
		return nil
	}

	seq := atomic.AddInt64(&d.readCursor, 1) % int64(len(d.buffer))
	return d.buffer[seq]
}

func (d *Disruptor) expandBuffer() {
	newBufferSize := len(d.buffer) * 2
	newBuffer := make([]interface{}, newBufferSize)

	for i := range d.buffer {
		newBuffer[i] = d.buffer[i]
	}

	d.buffer = newBuffer
}
