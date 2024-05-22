package order_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
)

func (u *OrderRepository) GetOrder(orderID string) (*order.Order, error) {
	stmt, err := u.dbConnection.Prepare(`SELECT orderId, userId, type, amount, price, status FROM orders WHERE orderId = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(orderID)

	var order order.Order
	err = row.Scan(&order.ID, &order.UserID, &order.Type, &order.Amount, &order.Price, &order.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("order %s not found", orderID))
		}
		return nil, err
	}

	return &order, nil
}

func (u *OrderRepository) GetPendingOrders(orderType int) ([]*order.Order, error) {
	stmt, err := u.dbConnection.Prepare("SELECT orderId, userId, type, amount, price, status FROM orders WHERE type = ? AND status = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(orderType, order.OrderStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var order order.Order

		err := rows.Scan(&order.ID, &order.UserID, &order.Type,
			&order.Amount, &order.Price, &order.Status)
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

func (u *OrderRepository) GetMemOrder(orderId string) (*order.Order, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	orderLocal, exists := u.orders[orderId]
	return orderLocal, exists
}
