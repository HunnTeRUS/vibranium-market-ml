package order_queue

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("SNAPSHOT_DIR", os.TempDir())
	code := m.Run()
	os.Unsetenv("SNAPSHOT_DIR")
	os.Exit(code)
}

func TestDisruptor(t *testing.T) {
	disruptor := NewDisruptor(2)
	order1 := &order.Order{ID: "order1", UserID: "user1", Price: 50.0, Amount: 100}
	order2 := &order.Order{ID: "order2", UserID: "user2", Price: 100.0, Amount: 200}

	t.Run("enqueue order successfully", func(t *testing.T) {
		err := disruptor.Enqueue(order1)
		assert.NoError(t, err)
		err = disruptor.Enqueue(order2)
		assert.NoError(t, err)
	})

	t.Run("dequeue order successfully", func(t *testing.T) {
		dequeuedOrder := disruptor.Dequeue()
		assert.Equal(t, order1, dequeuedOrder)
		dequeuedOrder = disruptor.Dequeue()
		assert.Equal(t, order2, dequeuedOrder)
	})

	t.Run("dequeue from empty queue", func(t *testing.T) {
		dequeuedOrder := disruptor.Dequeue()
		assert.Nil(t, dequeuedOrder)
	})
}
