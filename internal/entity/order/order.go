package order

import (
	"errors"
	"github.com/google/uuid"
)

const (
	OrderTypeBuy  = 1
	OrderTypeSell = 2

	OrderStatusPending   = "PENDING"
	OrderStatusCompleted = "COMPLETED"
	OrderStatusCanceled  = "CANCELED"
)

type Order struct {
	ID     string
	UserID string
	Type   int
	Amount int
	Price  float64
	Status string
}

type OrderRepositoryInterface interface {
	UpsertOrder(order *Order) error
	GetOrder(orderID string) (*Order, error)
}

type OrderQueueInterface interface {
	EnqueueOrder(order *Order) error
	DequeueOrder() (*Order, error)
}

func NewOrder(userID string, orderType int, amount int, price float64) (*Order, error) {
	if orderType != OrderTypeBuy && orderType != OrderTypeSell {
		return nil, errors.New("invalid order type")
	}
	if amount <= 0 {
		return nil, errors.New("invalid amount value")
	}
	if price <= 0 {
		return nil, errors.New("invalid price value")
	}

	return &Order{
		ID:     uuid.New().String(),
		UserID: userID,
		Type:   orderType,
		Amount: amount,
		Price:  price,
		Status: OrderStatusPending,
	}, nil
}
