package order_usecase

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"log"
	"os"
	"strconv"
)

type OrderInputDTO struct {
	UserID string  `json:"userId" binding:"required"`
	Type   int     `json:"type" binding:"required"`
	Amount int     `json:"amount" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

type OrderOutputDTO struct {
	ID     string  `json:"orderId"`
	UserID string  `json:"userId"`
	Type   int     `json:"type"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type OrderUsecaseInterface interface {
	ProcessOrder(orderInput *OrderInputDTO) (string, error)
	GetOrder(orderID string) (*OrderOutputDTO, error)
}

type OrderUsecase struct {
	orderRepositoryInterface  order.OrderRepositoryInterface
	queueInterface            order.OrderQueueInterface
	walletRepositoryInterface wallet.WalletRepositoryInterface
}

func NewOrderUsecase(
	orderRepositoryInterface order.OrderRepositoryInterface,
	walletRepositoryInterface wallet.WalletRepositoryInterface,
	queueInterface order.OrderQueueInterface,
) *OrderUsecase {
	orderUsecaseObj := &OrderUsecase{
		orderRepositoryInterface,
		queueInterface,
		walletRepositoryInterface}
	orderUsecaseObj.StartOrderProcessingWorker()
	return orderUsecaseObj
}

func (ou *OrderUsecase) StartOrderProcessingWorker() {
	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		numWorkers = 10
	}

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				orderUnprocessed, err := ou.queueInterface.DequeueOrder()
				if err != nil {
					continue
				}

				err = ou.ExecuteOrder(orderUnprocessed)
				if err != nil {
					log.Printf("Failed to process order: %v", err)
				}
			}
		}()
	}
}

func (ou *OrderUsecase) ProcessOrder(orderInput *OrderInputDTO) (string, error) {
	orderEntity, err := order.NewOrder(orderInput.UserID, orderInput.Type, orderInput.Amount, orderInput.Price)
	if err != nil {
		return "", err
	}

	err = ou.orderRepositoryInterface.UpsertOrder(orderEntity)
	if err != nil {
		return "", err
	}

	if err := ou.queueInterface.EnqueueOrder(orderEntity); err != nil {
		return "", err
	}

	return orderEntity.ID, nil
}

func (ou *OrderUsecase) ExecuteOrder(orderEntity *order.Order) error {
	switch orderEntity.Type {
	case order.OrderTypeBuy:
		walletEntity, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
		if err != nil {
			return err
		}
		if walletEntity.Balance < float64(orderEntity.Amount)*orderEntity.Price {
			orderEntity.Status = order.OrderStatusCanceled
			err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
			if err != nil {
				return err
			}

			return errors.New("insufficient balance")
		}

		walletEntity.Balance -= float64(orderEntity.Amount) * orderEntity.Price
		walletEntity.Vibranium += orderEntity.Amount

		err = ou.walletRepositoryInterface.UpdateWallet(walletEntity)
		if err != nil {
			return err
		}

	case order.OrderTypeSell:
		walletEntity, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
		if err != nil {
			return err
		}

		if walletEntity.Vibranium < orderEntity.Amount {
			orderEntity.Status = order.OrderStatusCanceled
			err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
			if err != nil {
				return err
			}

			return errors.New("insufficient vibranium")
		}

		walletEntity.Vibranium -= orderEntity.Amount
		walletEntity.Balance += float64(orderEntity.Amount) * orderEntity.Price
		err = ou.walletRepositoryInterface.UpdateWallet(walletEntity)
		if err != nil {
			return err
		}
	}

	orderEntity.Status = order.OrderStatusCompleted
	err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
	if err != nil {
		return err
	}

	return nil
}
