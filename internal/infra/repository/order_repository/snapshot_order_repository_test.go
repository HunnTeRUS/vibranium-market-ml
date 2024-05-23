package order_repository

import (
	"encoding/json"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSnapshot(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "orders_snapshot_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	os.Setenv("ORDERS_SNAPSHOT_FILE", tmpFile.Name())
	defer os.Unsetenv("ORDERS_SNAPSHOT_FILE")

	orders := map[string]*order.Order{
		"order1": {ID: "order1", UserID: "user1", Amount: 100, Price: 50.0},
		"order2": {ID: "order2", UserID: "user2", Amount: 200, Price: 100.0},
	}
	buyCache := map[int][]*order.Order{
		1: {orders["order1"]},
	}
	sellCache := map[int][]*order.Order{
		2: {orders["order2"]},
	}
	data, err := json.Marshal(struct {
		Orders    map[string]*order.Order `json:"orders"`
		BuyCache  map[int][]*order.Order  `json:"buy_cache"`
		SellCache map[int][]*order.Order  `json:"sell_cache"`
	}{
		Orders:    orders,
		BuyCache:  buyCache,
		SellCache: sellCache,
	})
	assert.NoError(t, err)

	_, err = tmpFile.Write(data)
	assert.NoError(t, err)
	tmpFile.Close()

	repo := NewOrderRepository()
	err = repo.LoadSnapshot()
	assert.NoError(t, err)

	repo.mu.Lock()
	defer repo.mu.Unlock()
	assert.Equal(t, orders, repo.orders)
	assert.Equal(t, buyCache, repo.buyCache)
	assert.Equal(t, sellCache, repo.sellCache)
}

func TestSaveSnapshot(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "orders_snapshot")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	snapshotFile := filepath.Join(tmpDir, "orders_snapshot.json")
	os.Setenv("ORDERS_SNAPSHOT_FILE", snapshotFile)
	defer os.Unsetenv("ORDERS_SNAPSHOT_FILE")

	repo := NewOrderRepository()
	repo.orders["order1"] = &order.Order{ID: "order1", UserID: "user1", Amount: 100, Price: 50.0}
	repo.orders["order2"] = &order.Order{ID: "order2", UserID: "user2", Amount: 200, Price: 100.0}
	repo.buyCache[1] = []*order.Order{repo.orders["order1"]}
	repo.sellCache[2] = []*order.Order{repo.orders["order2"]}

	err = repo.SaveSnapshot()
	assert.NoError(t, err)

	file, err := os.Open(snapshotFile)
	assert.NoError(t, err)
	defer file.Close()

	var data struct {
		Orders    map[string]*order.Order `json:"orders"`
		BuyCache  map[int][]*order.Order  `json:"buy_cache"`
		SellCache map[int][]*order.Order  `json:"sell_cache"`
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, repo.orders, data.Orders)
	assert.Equal(t, repo.buyCache, data.BuyCache)
	assert.Equal(t, repo.sellCache, data.SellCache)
}

func TestSaveSnapshot_EmptyOrders(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "orders_snapshot")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	snapshotFile := filepath.Join(tmpDir, "orders_snapshot.json")
	os.Setenv("ORDERS_SNAPSHOT_FILE", snapshotFile)
	defer os.Unsetenv("ORDERS_SNAPSHOT_FILE")

	repo := NewOrderRepository()

	err = repo.SaveSnapshot()
	assert.NoError(t, err)

	file, err := os.Open(snapshotFile)
	assert.NoError(t, err)
	defer file.Close()

	var data struct {
		Orders    map[string]*order.Order `json:"orders"`
		BuyCache  map[int][]*order.Order  `json:"buy_cache"`
		SellCache map[int][]*order.Order  `json:"sell_cache"`
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	assert.NoError(t, err)
	assert.Empty(t, data.Orders)
	assert.Empty(t, data.BuyCache)
	assert.Empty(t, data.SellCache)
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	os.Setenv("ORDERS_SNAPSHOT_FILE", "nonexistent_file.json")
	defer os.Unsetenv("ORDERS_SNAPSHOT_FILE")

	repo := NewOrderRepository()
	err := repo.LoadSnapshot()
	assert.Error(t, err)
}

func TestSaveSnapshot_DirectoryCreationError(t *testing.T) {
	snapshotFile := filepath.Join(string([]rune{0}), "orders_snapshot.json") // Invalid directory path
	os.Setenv("ORDERS_SNAPSHOT_FILE", snapshotFile)
	defer os.Unsetenv("ORDERS_SNAPSHOT_FILE")

	repo := NewOrderRepository()
	repo.orders["order1"] = &order.Order{ID: "order1", UserID: "user1", Amount: 100, Price: 50.0}

	err := repo.SaveSnapshot()
	assert.Error(t, err)
}
