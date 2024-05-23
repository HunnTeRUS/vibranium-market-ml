package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"sync"
)

type OrderRepository struct {
	mu        sync.RWMutex
	taskQueue chan *order.Order

	orders    map[string]*order.Order
	buyCache  map[int][]*order.Order
	sellCache map[int][]*order.Order
}

func NewOrderRepository() *OrderRepository {
	orderRepo := &OrderRepository{
		orders:    make(map[string]*order.Order),
		taskQueue: make(chan *order.Order, 20000),
		buyCache:  make(map[int][]*order.Order),
		sellCache: make(map[int][]*order.Order),
	}

	return orderRepo
}

func (u *OrderRepository) UpsertOrder(order *order.Order) {
	u.UpsertLocalOrder(order)
}

func (u *OrderRepository) UpsertLocalOrder(orderEntity *order.Order) {
	u.orders[orderEntity.ID] = orderEntity

	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)

	if orderEntity.Status != order.OrderStatusPending {
		u.removeOrderFromCache(orderEntity)
	} else {
		if orderEntity.Type == order.OrderTypeBuy {
			u.addOrderToCache(&u.buyCache, priceBucket, orderEntity)
		} else {
			u.addOrderToCache(&u.sellCache, priceBucket, orderEntity)
		}
	}
}
