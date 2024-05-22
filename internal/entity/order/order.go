package order

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/metrics"
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
	Symbol string
	Status string
}

type OrderRepositoryInterface interface {
	UpsertOrder(order *Order)
	GetOrder(orderID string) (*Order, error)
	GetMemOrder(orderId string) (*Order, bool)
	GetPendingOrders(symbol string, orderType int) ([]*Order, error)
}

type OrderQueueInterface interface {
	EnqueueOrder(order *Order) error
	DequeueOrder() (*Order, error)
}

func (o *Order) CompleteOrder(orderRepositoryInterface OrderRepositoryInterface) error {
	o.Status = OrderStatusCompleted

	metrics.OrderProcessed.Inc()
	orderRepositoryInterface.UpsertOrder(o)

	return nil
}

func (o *Order) CancelOrder(orderRepositoryInterface OrderRepositoryInterface, reason string) error {
	o.Status = OrderStatusCanceled

	metrics.OrderCanceled.Inc()
	orderRepositoryInterface.UpsertOrder(o)

	return errors.New(reason)
}

func NewOrder(userID string, orderType int, amount int, price float64, symbol string) (*Order, error) {
	if orderType != OrderTypeBuy && orderType != OrderTypeSell {
		return nil, errors.New("invalid order type")
	}
	if amount <= 0 {
		return nil, errors.New("invalid amount value")
	}
	if price <= 0 {
		return nil, errors.New("invalid price value")
	}
	if len(symbol) != 3 {
		return nil, errors.New("invalid stock name")
	}

	return &Order{
		ID:     uuid.New().String(),
		UserID: userID,
		Type:   orderType,
		Amount: amount,
		Price:  price,
		Symbol: symbol,
		Status: OrderStatusPending,
	}, nil
}
