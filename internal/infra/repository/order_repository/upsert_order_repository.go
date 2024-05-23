package order_repository

import (
	"encoding/json"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"os"
	"sync"
)

type OrderRepository struct {
	mu        sync.RWMutex
	taskQueue chan *order.Order

	orders    map[string]*order.Order
	buyCache  *sync.Map
	sellCache *sync.Map
}

func NewOrderRepository() *OrderRepository {
	orderRepo := &OrderRepository{
		orders:    make(map[string]*order.Order),
		taskQueue: make(chan *order.Order, 20000),
		buyCache:  &sync.Map{},
		sellCache: &sync.Map{},
	}

	return orderRepo
}

func (u *OrderRepository) LoadSnapshot() error {
	file, err := os.Open(os.Getenv("ORDERS_SNAPSHOT_FILE"))
	if err != nil {
		logger.Warn("environment variable ORDERS_SNAPSHOT_FILE not set")
		return err
	}
	defer file.Close()

	data := struct {
		Orders    map[string]*order.Order   `json:"orders"`
		BuyCache  map[string][]*order.Order `json:"buy_cache"`
		SellCache map[string][]*order.Order `json:"sell_cache"`
	}{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		logger.Warn(err.Error())
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.orders = data.Orders

	for key, orders := range data.BuyCache {
		u.buyCache.Store(key, orders)
	}

	for key, orders := range data.SellCache {
		u.sellCache.Store(key, orders)
	}

	return nil
}

func (u *OrderRepository) SaveSnapshot() error {
	u.mu.RLock()
	defer u.mu.RUnlock()

	buyCache := make(map[string][]*order.Order)
	u.buyCache.Range(func(key, value interface{}) bool {
		buyCache[key.(string)] = value.([]*order.Order)
		return true
	})

	sellCache := make(map[string][]*order.Order)
	u.sellCache.Range(func(key, value interface{}) bool {
		sellCache[key.(string)] = value.([]*order.Order)
		return true
	})

	data := struct {
		Orders    map[string]*order.Order   `json:"orders,omitempty"`
		BuyCache  map[string][]*order.Order `json:"buy_cache,omitempty"`
		SellCache map[string][]*order.Order `json:"sell_cache,omitempty"`
	}{
		Orders:    u.orders,
		BuyCache:  buyCache,
		SellCache: sellCache,
	}

	file, err := os.Create(os.Getenv("ORDERS_SNAPSHOT_FILE"))
	if err != nil {
		logger.Warn(err.Error())
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		logger.Warn(err.Error())
		return err
	}

	return nil
}

func (u *OrderRepository) UpsertOrder(order *order.Order) {
	u.UpsertLocalOrder(order)
}

func (u *OrderRepository) UpsertLocalOrder(orderEntity *order.Order) {
	u.mu.Lock()
	u.orders[orderEntity.ID] = orderEntity
	u.mu.Unlock()
}
