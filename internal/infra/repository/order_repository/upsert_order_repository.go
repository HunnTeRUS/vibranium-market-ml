package order_repository

import (
	"database/sql"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"sync"
)

type OrderRepository struct {
	sync.Mutex
	dbConnection *sql.DB

	orders map[string]*order.Order
}

func NewOrderRepository(dbConnection *sql.DB) *OrderRepository {
	return &OrderRepository{dbConnection: dbConnection, orders: make(map[string]*order.Order)}
}

func (u *OrderRepository) UpsertOrder(order *order.Order) error {
	u.UpsertLocalOrder(order)

	go func() {
		if _, err := u.GetOrder(order.ID); err != nil {
			stmt, err := u.dbConnection.Prepare("INSERT INTO orders " +
				"(id, user_id, type, amount, price, status, symbol) VALUES (?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				logger.Error("error trying to prepare insert statement", err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(order.ID, order.UserID, order.Type, order.Amount, order.Price, order.Status, order.Symbol)
			if err != nil {
				logger.Error("error trying to insert order", err)
			}
		} else {
			stmt, err := u.dbConnection.Prepare("UPDATE orders SET status = ? WHERE id = ?")
			if err != nil {
				logger.Error("error trying to prepare update order", err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(order.Status, order.ID)
			if err != nil {
				logger.Error("error trying to update order", err)
			}
		}
	}()

	return nil
}

func (u *OrderRepository) UpsertLocalOrder(orderEntity *order.Order) {
	u.Lock()
	defer u.Unlock()
	u.orders[orderEntity.ID] = orderEntity
}
