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
	Status string

	SellValueRemaining int
	CorrelationId      string
}

type OrderRepositoryInterface interface {
	UpsertOrder(order *Order)
	GetMemOrder(orderId string) (*Order, bool)
	GetBuyingMatchingOrder(orderEntity *Order) (*Order, error)
	GetSellingMatchingOrder(orderEntity *Order) (*Order, error)

	LoadSnapshot() error
	SaveSnapshot() error
}

type OrderQueueInterface interface {
	EnqueueOrder(order *Order) error
	DequeueOrder() (*Order, error)
}

func (o *Order) CompleteOrder(orderRepositoryInterface OrderRepositoryInterface) {
	o.Status = OrderStatusCompleted

	metrics.OrderProcessed.Inc()
	orderRepositoryInterface.UpsertOrder(o)
}

func (o *Order) CancelOrder(orderRepositoryInterface OrderRepositoryInterface, reason string) error {
	o.Status = OrderStatusCanceled

	metrics.OrderCanceled.Inc()
	orderRepositoryInterface.UpsertOrder(o)

	return errors.New(reason)
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

	orderValue := &Order{
		ID:     uuid.New().String(),
		UserID: userID,
		Type:   orderType,
		Amount: amount,
		Price:  price,
		Status: OrderStatusPending,
	}

	if orderType == OrderTypeSell {
		orderValue.SellValueRemaining = amount
	}

	return orderValue, nil
}
