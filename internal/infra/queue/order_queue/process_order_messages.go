package order_queue

import (
	"encoding/json"
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/metrics"
)

type OrderQueue struct {
	disruptor *Disruptor
}

func NewOrderQueue(initialSize int) *OrderQueue {
	return &OrderQueue{
		disruptor: NewDisruptor(initialSize),
	}
}

func (q *OrderQueue) EnqueueOrder(order *order.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		metrics.ProcessingErrors.Inc()
		return err
	}
	err = q.disruptor.Enqueue(orderJSON)
	if err == ErrBufferFull {
		metrics.ProcessingErrors.Inc()
		metrics.BufferFull.Inc()
		logger.Warn("Buffer is full, expanding buffer and retrying...")
		err = q.disruptor.Enqueue(orderJSON)
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

	var orderEntity order.Order
	t, ok := message.([]byte)
	if !ok {
		return nil, errors.New("invalid object type inside the queue")
	}

	err := json.Unmarshal(t, &orderEntity)
	if err != nil {
		metrics.ProcessingErrors.Inc()
		return nil, err
	}

	metrics.OrdersDequeued.Inc()

	return &orderEntity, nil
}

func (q *OrderQueue) SaveSnapshot() error {
	return q.disruptor.SaveSnapshotToS3()
}

func (q *OrderQueue) LoadSnapshot() error {
	return q.disruptor.LoadSnapshotFromS3()
}
