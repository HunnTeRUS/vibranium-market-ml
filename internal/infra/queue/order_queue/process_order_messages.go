package order_queue

import (
	"errors"
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

func (q *OrderQueue) EnqueueOrder(order *order.Order) error {
	metrics.OrdersEnqueued.Inc()

	return q.disruptor.Enqueue(order)
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
