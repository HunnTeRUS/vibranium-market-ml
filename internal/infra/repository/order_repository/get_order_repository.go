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

func (u *OrderRepository) GetSellingMatchingOrder(orderEntity *order.Order) (*order.Order, error) {
	if orderEntity.Type != order.OrderTypeBuy {
		return nil, nil
	}

	var matchingOrder *order.Order
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

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

func (u *OrderRepository) addOrderToCache(cache *map[int][]*order.Order, priceBucket int, orderEntity *order.Order) {
	(*cache)[priceBucket] = append((*cache)[priceBucket], orderEntity)
}

func (u *OrderRepository) removeOrderFromCache(orderEntity *order.Order) {
	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

	if orderEntity.Type == order.OrderTypeBuy {
		u.buyCache[priceBucket] = removeOrderFromSlice(u.buyCache[priceBucket], orderEntity)
	} else {
		u.sellCache[priceBucket] = removeOrderFromSlice(u.sellCache[priceBucket], orderEntity)
	}
}

func removeOrderFromSlice(slice []*order.Order, orderEntity *order.Order) []*order.Order {
	for i, ord := range slice {
		if ord.ID == orderEntity.ID {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (u *OrderRepository) GetMemOrder(orderId string) (*order.Order, bool) {
	orderLocal, exists := u.orders[orderId]
	return orderLocal, exists
}
