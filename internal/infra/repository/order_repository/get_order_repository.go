package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
)

const priceBucketSize = 10.0
const amountBucketSize = 10

func calculateBucket(value, bucketSize float64) int {
	return int(value / bucketSize)
}

func (u *OrderRepository) GetBuyingMatchingOrder(orderEntity *order.Order) (*order.Order, error) {
	if orderEntity.Type != order.OrderTypeSell {
		return nil, nil
	}

	var matchingOrder *order.Order
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

	u.mu.RLock()
	defer u.mu.RUnlock()

	if buyOrders, exists := u.buyCache[priceBucket]; exists {
		for _, buyOrder := range buyOrders {
			if buyOrder.Price == orderEntity.Price && buyOrder.Amount <= orderEntity.SellValueRemaining && buyOrder.UserID != orderEntity.UserID {
				matchingOrder = buyOrder
				break
			}
		}
	}

	return matchingOrder, nil
}

// GetSellingMatchingOrder finds the first sell order that matches the criteria of the buy order
func (u *OrderRepository) GetSellingMatchingOrder(orderEntity *order.Order) (*order.Order, error) {
	if orderEntity.Type != order.OrderTypeBuy {
		return nil, nil
	}

	var matchingOrder *order.Order
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

	u.mu.RLock()
	defer u.mu.RUnlock()

	if sellOrders, exists := u.sellCache[priceBucket]; exists {
		for _, sellOrder := range sellOrders {
			if sellOrder.Price == orderEntity.Price && sellOrder.SellValueRemaining >= orderEntity.Amount && sellOrder.UserID != orderEntity.UserID {
				matchingOrder = sellOrder
				break
			}
		}
	}

	return matchingOrder, nil
}

func (u *OrderRepository) addOrderToCache(cache *map[int][]*order.Order, priceBucket int, o *order.Order) {
	if orders, exists := (*cache)[priceBucket]; exists {
		(*cache)[priceBucket] = append(orders, o)
	} else {
		(*cache)[priceBucket] = []*order.Order{o}
	}
}

func (u *OrderRepository) removeOrderFromCache(orderEntity *order.Order) {
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

	var cache *map[int][]*order.Order
	if orderEntity.Type == order.OrderTypeBuy {
		cache = &u.buyCache
	} else {
		cache = &u.sellCache
	}

	if orders, exists := (*cache)[priceBucket]; exists {
		for i, o := range orders {
			if o.ID == orderEntity.ID {
				(*cache)[priceBucket] = append(orders[:i], orders[i+1:]...)
				break
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
