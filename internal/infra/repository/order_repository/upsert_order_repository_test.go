package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpsertOrder(t *testing.T) {
	repo := NewOrderRepository(10)

	t.Run("insert new buy order", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100, Status: order.OrderStatusPending}
		repo.UpsertOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		createdOrder, exists := repo.orders["order1"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, createdOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.Contains(t, repo.buyCache[priceBucket], orderEntity)
	})

	t.Run("insert new sell order", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order2", UserID: "user2", Type: order.OrderTypeSell, Price: 100.0, Amount: 200, Status: order.OrderStatusPending}
		repo.UpsertOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		createdOrder, exists := repo.orders["order2"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, createdOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.Contains(t, repo.sellCache[priceBucket], orderEntity)
	})

	t.Run("update existing order", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 150, Status: order.OrderStatusPending}
		repo.UpsertOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		updatedOrder, exists := repo.orders["order1"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, updatedOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.Contains(t, repo.buyCache[priceBucket], orderEntity)
	})

	t.Run("remove non-pending order from cache", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 150, Status: order.OrderStatusCompleted}
		repo.UpsertOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		updatedOrder, exists := repo.orders["order1"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, updatedOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.NotContains(t, repo.buyCache[priceBucket], orderEntity)
	})
}

func TestUpsertLocalOrder(t *testing.T) {
	repo := NewOrderRepository(10)

	t.Run("insert new buy order locally", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100, Status: order.OrderStatusPending}
		repo.UpsertLocalOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		createdOrder, exists := repo.orders["order1"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, createdOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.Contains(t, repo.buyCache[priceBucket], orderEntity)
	})

	t.Run("insert new sell order locally", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order2", UserID: "user2", Type: order.OrderTypeSell, Price: 100.0, Amount: 200, Status: order.OrderStatusPending}
		repo.UpsertLocalOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		createdOrder, exists := repo.orders["order2"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, createdOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.Contains(t, repo.sellCache[priceBucket], orderEntity)
	})

	t.Run("update existing order locally", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 150, Status: order.OrderStatusPending}
		repo.UpsertLocalOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		updatedOrder, exists := repo.orders["order1"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, updatedOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.Contains(t, repo.buyCache[priceBucket], orderEntity)
	})

	t.Run("remove non-pending order from cache locally", func(t *testing.T) {
		orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 150, Status: order.OrderStatusCompleted}
		repo.UpsertLocalOrder(orderEntity)

		time.Sleep(10 * time.Millisecond)

		updatedOrder, exists := repo.orders["order1"]

		assert.True(t, exists)
		assert.Equal(t, orderEntity, updatedOrder)

		priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
		assert.NotContains(t, repo.buyCache[priceBucket], orderEntity)
	})
}
