package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
)

type OrderRepository struct {
	orderChan   chan *order.Order
	cacheUpdate chan *cacheUpdate

	orders    map[string]*order.Order
	buyCache  map[int][]*order.Order
	sellCache map[int][]*order.Order
}

type cacheUpdate struct {
	orderEntity *order.Order
	priceBucket int
}

func NewOrderRepository(bufferSize int) *OrderRepository {
	repo := &OrderRepository{
		orderChan:   make(chan *order.Order, bufferSize),
		cacheUpdate: make(chan *cacheUpdate, bufferSize),
		orders:      make(map[string]*order.Order, bufferSize),
		buyCache:    make(map[int][]*order.Order, bufferSize),
		sellCache:   make(map[int][]*order.Order, bufferSize),
	}

	go repo.processOrders()
	go repo.processCacheUpdates()

	return repo
}

func (u *OrderRepository) processOrders() {
	for orderEntity := range u.orderChan {
		u.orders[orderEntity.ID] = orderEntity
	}
}

func (u *OrderRepository) processCacheUpdates() {
	if u.buyCache == nil {
		u.buyCache = make(map[int][]*order.Order)
	}

	if u.sellCache == nil {
		u.sellCache = make(map[int][]*order.Order)
	}

	for update := range u.cacheUpdate {
		orderEntity := update.orderEntity
		priceBucket := update.priceBucket

		if orderEntity.Status != order.OrderStatusPending {
			u.removeOrderFromCache(orderEntity)
		} else {
			if u.buyCache == nil {
				u.buyCache = make(map[int][]*order.Order)
			}

			if u.sellCache == nil {
				u.sellCache = make(map[int][]*order.Order)
			}

			if orderEntity.Type == order.OrderTypeBuy {
				u.addOrderToCache(&u.buyCache, priceBucket, orderEntity)
			} else {
				u.addOrderToCache(&u.sellCache, priceBucket, orderEntity)
			}
		}
	}
}

func (u *OrderRepository) UpsertOrder(order *order.Order) {
	u.UpsertLocalOrder(order)
}

func (u *OrderRepository) UpsertLocalOrder(orderEntity *order.Order) {
	u.orderChan <- orderEntity

	priceBucket := calculateBucket(orderEntity.Price, priceBucketSize)
	u.cacheUpdate <- &cacheUpdate{orderEntity, priceBucket}
}
