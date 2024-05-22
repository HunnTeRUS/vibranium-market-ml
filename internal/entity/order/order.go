package order

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
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
	UpsertOrder(order *Order) error
	GetOrder(orderID string) (*Order, error)
	GetPendingOrders(symbol string, orderType int) ([]*Order, error)
}

type OrderQueueInterface interface {
	EnqueueOrder(order *Order) error
	DequeueOrder() (*Order, error)
}

func (o *Order) CompleteOrder(orderRepositoryInterface OrderRepositoryInterface) error {
	o.Status = OrderStatusCompleted

	metrics.OrderProcessed.Inc()
	err := orderRepositoryInterface.UpsertOrder(o)
	if err != nil {
		metrics.ProcessingErrors.Inc()
		logger.Error("action=ExecuteOrder, "+
			"message=error calling UpsertOrder repository for completing order", err)
		return err
	}
	return nil
}

func (o *Order) CancelOrder(orderRepositoryInterface OrderRepositoryInterface, reason string) error {
	o.Status = OrderStatusCanceled

	metrics.OrderCanceled.Inc()
	err := orderRepositoryInterface.UpsertOrder(o)
	if err != nil {
		logger.Error("action=ExecuteOrder, "+
			"message=error calling UpsertOrder repository for cancelling order", err)
		return err
	}
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

	return &Order{
		ID:     uuid.New().String(),
		UserID: userID,
		Type:   orderType,
		Amount: amount,
		Price:  price,
		Status: OrderStatusPending,
	}, nil
}
