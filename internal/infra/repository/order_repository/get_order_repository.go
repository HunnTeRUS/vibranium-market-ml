package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"sync"
)

const priceBucketSize = 10.0
const amountBucketSize = 10

func calculateBucket(value, bucketSize float64) int {
	return int(value / bucketSize)
}

func (u *OrderRepository) GetFirstMatchingOrder(orderEntity *order.Order) (*order.Order, error) {
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
	amountBucket := calculateBucket(float64(orderEntity.Amount), amountBucketSize)

	var cache *sync.Map
	if orderEntity.Type == order.OrderTypeBuy {
		cache = u.buyCache
	} else {
		cache = u.sellCache
	}

	cacheKey := (priceBucket << 32) | amountBucket
	if ordersInterface, ok := cache.Load(cacheKey); ok {
		if orders, ok := ordersInterface.([]*order.Order); ok && len(orders) > 0 {
			for _, orderValue := range orders {
				if orderValue.UserID == orderEntity.UserID {
					continue
				}
				return orderValue, nil
			}
		}
	}

	return nil, nil
}

func (u *OrderRepository) addOrderToCache(cache *sync.Map, priceBucket, amountBucket int, o *order.Order) {
	cacheKey := (priceBucket << 32) | amountBucket
	u.mu.Lock()
	defer u.mu.Unlock()
	if ordersInterface, ok := cache.Load(cacheKey); ok {
		if orders, ok := ordersInterface.([]*order.Order); ok {
			cache.Store(cacheKey, append(orders, o))
			return
		}
	}
	cache.Store(cacheKey, []*order.Order{o})
}

func (u *OrderRepository) removeOrderFromCache(orderType int, price float64, amount int, orderID string) {
	priceBucket := calculateBucket(price, priceBucketSize)
	amountBucket := calculateBucket(float64(amount), amountBucketSize)

	var cache *sync.Map
	if orderType == 1 {
		cache = u.buyCache
	} else {
		cache = u.sellCache
	}

	cacheKey := (priceBucket << 32) | amountBucket
	u.mu.Lock()
	defer u.mu.Unlock()
	if ordersInterface, ok := cache.Load(cacheKey); ok {
		if orders, ok := ordersInterface.([]*order.Order); ok {
			for i, o := range orders {
				if o.ID == orderID {
					cache.Store(cacheKey, append(orders[:i], orders[i+1:]...))
					return
				}
			}
		}
	}
}

func (u *OrderRepository) GetMemOrder(orderId string) (*order.Order, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	orderLocal, exists := u.orders[orderId]
	return orderLocal, exists
}
