package order_queue

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderQueue(t *testing.T) {
	orderQueue := NewOrderQueue(2)
	order1 := &order.Order{ID: "order1", UserID: "user1", Price: 50.0, Amount: 100}
	order2 := &order.Order{ID: "order2", UserID: "user2", Price: 100.0, Amount: 200}

	t.Run("enqueue order successfully", func(t *testing.T) {
		err := orderQueue.EnqueueOrder(order1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(orderQueue.disruptor.queue))

		err = orderQueue.EnqueueOrder(order2)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(orderQueue.disruptor.queue))

		orderQueue.DequeueOrder()
		orderQueue.DequeueOrder()
	})

	t.Run("dequeue from empty queue", func(t *testing.T) {
		err := orderQueue.EnqueueOrder(nil)

		_, err = orderQueue.DequeueOrder()
		assert.Equal(t, "invalid nil object inside the queue", err.Error())
	})
}
