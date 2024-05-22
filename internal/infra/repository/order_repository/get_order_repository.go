package order_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
)

func (u *OrderRepository) GetOrder(orderID string) (*order.Order, error) {
	stmt, err := u.dbConnection.Prepare("SELECT * FROM orders WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(orderID)

	var order order.Order
	err = row.Scan(&order.ID, &order.UserID, &order.Type, &order.Amount, &order.Price, &order.Status, &order.Symbol)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(fmt.Sprintf("order %s not found", orderID))
			return nil, errors.New(fmt.Sprintf("order %s not found", orderID))
		}
		return nil, err
	}

	u.UpsertLocalOrder(&order)

	return &order, nil
}

func (u *OrderRepository) GetPendingOrders(symbol string, orderType int) ([]*Order, error) {
	stmt, err := u.dbConnection.Prepare("SELECT * FROM orders WHERE symbol = ? AND type = ? AND status = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(symbol, orderType, order.OrderStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var order order.Order

		err := rows.Scan(&order.ID, &order.UserID, &order.Symbol, &order.Type,
			&order.Amount, &order.Price, &order.Status, &order.Symbol)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (u *OrderRepository) GetLocalOrder(orderId string) (*order.Order, bool) {
	u.Lock()
	defer u.Unlock()
	orderLocal, exists := u.orders[orderId]
	return orderLocal, exists
}
