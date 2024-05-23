package order_queue

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"sync"
)

var ErrBufferFull = errors.New("buffer is full")

type Disruptor struct {
	queue       chan *order.Order
	writeCursor int64
	readCursor  int64
	mu          sync.Mutex
}

func NewDisruptor(bufferSize int) *Disruptor {
	d := &Disruptor{
		queue: make(chan *order.Order, bufferSize),
	}

	err := d.LoadSnapshotFromFile()
	if err != nil {
		logger.Warn("No previous snapshot found or failed to load")
	}

	return d
}

func (d *Disruptor) Enqueue(order *order.Order) error {
	select {
	case d.queue <- order:
		return nil
	default:
		return ErrBufferFull
	}
}

func (d *Disruptor) Dequeue() *order.Order {
	select {
	case order := <-d.queue:
		return order
	default:
		return nil
	}
}

func (q *OrderQueue) expandBuffer() {
	newBufferSize := q.bufferSize * 2
	newQueue := make(chan *order.Order, newBufferSize)

	for {
		select {
		case order := <-q.disruptor.queue:
			newQueue <- order
		default:
			q.disruptor.queue = newQueue
			q.bufferSize = newBufferSize
			return
		}
	}
}
