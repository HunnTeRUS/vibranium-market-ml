package order_repository

import (
	"database/sql"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"log"
	"sync"
)

type OrderRepository struct {
	mu           sync.RWMutex
	taskQueue    chan *order.Order
	dbConnection *sql.DB

	orders map[string]*order.Order
}

func NewOrderRepository(dbConnection *sql.DB) *OrderRepository {
	orderRepo := &OrderRepository{
		dbConnection: dbConnection,
		orders:       make(map[string]*order.Order),
		taskQueue:    make(chan *order.Order, 1000),
	}

	go orderRepo.worker()

	return orderRepo
}

func (u *OrderRepository) worker() {
	for {
		select {
		case order := <-u.taskQueue:
			u.upsertOrderInDB(order)
		}
	}
}

func (u *OrderRepository) upsertOrderInDB(order *order.Order) {
	if orderStored, _ := u.GetOrder(order.ID); orderStored == nil {
		stmt, err := u.dbConnection.Prepare("INSERT INTO orders (orderId, userId, type, amount, price, status) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Println("error trying to prepare insert statement", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(order.ID, order.UserID, order.Type, order.Amount, order.Price, order.Status)
		if err != nil {
			log.Println("error trying to insert order", err)
			return
		}
	} else {
		stmt, err := u.dbConnection.Prepare("UPDATE orders SET status = ? WHERE orderId = ?")
		if err != nil {
			log.Println("error trying to prepare update order", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(order.Status, order.ID)
		if err != nil {
			log.Println("error trying to update order", err)
			return
		}
	}
}

func (u *OrderRepository) UpsertOrder(order *order.Order) {
	u.UpsertLocalOrder(order)

	u.taskQueue <- order
}

func (u *OrderRepository) UpsertLocalOrder(orderEntity *order.Order) {
	u.mu.Lock()
	u.orders[orderEntity.ID] = orderEntity
	u.mu.Unlock()
}
