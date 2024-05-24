package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateBucket(t *testing.T) {
	tests := []struct {
		value      float64
		bucketSize float64
		expected   int
	}{
		{value: 10, bucketSize: 5, expected: 2},
		{value: 15, bucketSize: 5, expected: 3},
		{value: 20, bucketSize: 10, expected: 2},
		{value: 25, bucketSize: 10, expected: 2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := calculateBucket(tt.value, tt.bucketSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBuyingMatchingOrder(t *testing.T) {
	repo := NewOrderRepository(1010)

	orderSell := &order.Order{ID: "orderSell", UserID: "user1", Type: order.OrderTypeSell, Price: 50.0, Amount: 100, SellValueRemaining: 100}
	orderBuy := &order.Order{ID: "orderBuy", UserID: "user2", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100}

	repo.addOrderToCache(&repo.buyCache, calculateBucket(orderBuy.Price, priceBucketSize), orderBuy)

	t.Run("matching buy order found", func(t *testing.T) {
		matchingOrder, err := repo.GetBuyingMatchingOrder(orderSell)
		assert.NoError(t, err)
		assert.Equal(t, orderBuy, matchingOrder)
	})

	t.Run("no matching buy order found", func(t *testing.T) {
		orderSellNoMatch := &order.Order{ID: "orderSellNoMatch", UserID: "user1", Type: order.OrderTypeSell, Price: 100.0, Amount: 100, SellValueRemaining: 100}
		matchingOrder, err := repo.GetBuyingMatchingOrder(orderSellNoMatch)
		assert.NoError(t, err)
		assert.Nil(t, matchingOrder)
	})
}

func TestGetSellingMatchingOrder(t *testing.T) {
	repo := NewOrderRepository(10)

	orderBuy := &order.Order{ID: "orderBuy", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100}
	orderSell := &order.Order{ID: "orderSell", UserID: "user2", Type: order.OrderTypeSell, Price: 50.0, Amount: 100, SellValueRemaining: 100}

	repo.addOrderToCache(&repo.sellCache, calculateBucket(orderSell.Price, priceBucketSize), orderSell)

	t.Run("matching sell order found", func(t *testing.T) {
		matchingOrder, err := repo.GetSellingMatchingOrder(orderBuy)
		assert.NoError(t, err)
		assert.Equal(t, orderSell, matchingOrder)
	})

	t.Run("no matching sell order found", func(t *testing.T) {
		orderBuyNoMatch := &order.Order{ID: "orderBuyNoMatch", UserID: "user1", Type: order.OrderTypeBuy, Price: 100.0, Amount: 100}
		matchingOrder, err := repo.GetSellingMatchingOrder(orderBuyNoMatch)
		assert.NoError(t, err)
		assert.Nil(t, matchingOrder)
	})
}

func TestAddOrderToCache(t *testing.T) {
	repo := NewOrderRepository(10)

	orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100}
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

	repo.addOrderToCache(&repo.buyCache, priceBucket, orderEntity)

	assert.Contains(t, repo.buyCache[priceBucket], orderEntity)
}

func TestRemoveOrderFromCache(t *testing.T) {
	repo := NewOrderRepository(10)

	orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100}
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
	repo.addOrderToCache(&repo.buyCache, priceBucket, orderEntity)

	repo.removeOrderFromCache(orderEntity)

	assert.NotContains(t, repo.buyCache[priceBucket], orderEntity)
}

func TestGetMemOrder(t *testing.T) {
	repo := NewOrderRepository(10)

	orderEntity := &order.Order{ID: "order1", UserID: "user1", Type: order.OrderTypeBuy, Price: 50.0, Amount: 100}
	repo.UpsertLocalOrder(orderEntity)

	t.Run("order exists", func(t *testing.T) {
		result, exists := repo.GetMemOrder(orderEntity.ID)
		assert.True(t, exists)
		assert.Equal(t, orderEntity, result)
	})

	t.Run("order does not exist", func(t *testing.T) {
		result, exists := repo.GetMemOrder("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, result)
	})
}
