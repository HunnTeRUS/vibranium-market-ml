package order_queue

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/metrics"
	"sync"
)

type OrderQueue struct {
	queue      chan *order.Order
	disruptor  *Disruptor
	mu         sync.Mutex
	bufferSize int
}

func NewOrderQueue(initialSize int) *OrderQueue {
	return &OrderQueue{
		disruptor: NewDisruptor(initialSize),
	}
}

func (q *OrderQueue) QueueLength() int {
	return len(q.disruptor.queue)
}

func (q *OrderQueue) EnqueueOrder(order *order.Order) error {
	err := q.disruptor.Enqueue(order)
	if err == ErrBufferFull {
		metrics.ProcessingErrors.Inc()
		metrics.BufferFull.Inc()
		logger.Warn("Buffer is full, expanding buffer and retrying...")

		q.mu.Lock()
		q.expandBuffer()
		q.mu.Unlock()

		err = q.disruptor.Enqueue(order)
		if err != nil {
			return err
		}
	}
	metrics.OrdersEnqueued.Inc()
	return nil
}

func (q *OrderQueue) DequeueOrder() (*order.Order, error) {
	message := q.disruptor.Dequeue()
	if message == nil {
		return nil, errors.New("invalid nil object inside the queue")
	}

	metrics.OrdersDequeued.Inc()

	return message, nil
}

func (q *OrderQueue) SaveSnapshot() error {
	return q.disruptor.SaveSnapshotToFile()
}

func (q *OrderQueue) LoadSnapshot() error {
	return q.disruptor.LoadSnapshotFromFile()
}
