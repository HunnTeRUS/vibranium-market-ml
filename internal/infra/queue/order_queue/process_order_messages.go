package order_queue

import (
	"encoding/json"
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

func (q *OrderQueue) DequeueOrder() ([]*order.Order, error) {
	var orders []*order.Order

	for i := 0; i < 10; i++ {
		message := q.disruptor.Dequeue()
		if message == nil {
			break
		}

		var orderEntity order.Order
		err := json.Unmarshal(message.([]byte), &orderEntity)
		if err != nil {
			metrics.ProcessingErrors.Inc()
			return nil, err
		}
		orders = append(orders, &orderEntity)
	}
	metrics.OrdersDequeued.Add(float64(len(orders)))
	return orders, nil
}

func (q *OrderQueue) SaveSnapshot() error {
	return q.disruptor.SaveSnapshotToS3()
}

func (q *OrderQueue) LoadSnapshot() error {
	return q.disruptor.LoadSnapshotFromS3()
}
