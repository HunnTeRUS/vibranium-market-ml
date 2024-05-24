package order_repository

import (
	"encoding/json"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"os"
	"path/filepath"
)

func (u *OrderRepository) LoadSnapshot() error {
	file, err := os.Open(os.Getenv("ORDERS_SNAPSHOT_FILE"))
	if err != nil {
		logger.Warn("environment variable ORDERS_SNAPSHOT_FILE not set")
		return err
	}
	defer file.Close()

	data := struct {
		Orders    map[string]*order.Order `json:"orders"`
		BuyCache  map[int][]*order.Order  `json:"buy_cache"`
		SellCache map[int][]*order.Order  `json:"sell_cache"`
	}{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		logger.Warn(err.Error())
		return err
	}

	if data.Orders == nil {
		u.orders = make(map[string]*order.Order, 30000)
	} else {
		u.orders = data.Orders
	}

	u.buyCache = data.BuyCache
	u.sellCache = data.SellCache

	return nil
}

func (u *OrderRepository) SaveSnapshot() error {
	ordersSnapshotFile := os.Getenv("ORDERS_SNAPSHOT_FILE")

	snapshotDir := filepath.Dir(ordersSnapshotFile)
	err := os.MkdirAll(snapshotDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	data := struct {
		Orders    map[string]*order.Order `json:"orders,omitempty"`
		BuyCache  map[int][]*order.Order  `json:"buy_cache,omitempty"`
		SellCache map[int][]*order.Order  `json:"sell_cache,omitempty"`
	}{
		Orders:    u.orders,
		BuyCache:  u.buyCache,
		SellCache: u.sellCache,
	}

	file, err := os.Create(ordersSnapshotFile)
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
