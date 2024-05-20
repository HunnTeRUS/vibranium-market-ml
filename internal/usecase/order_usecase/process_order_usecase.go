package order_usecase

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"log"
	"os"
	"strconv"
)

const (
	NumWorkers = 10
)

type OrderInputDTO struct {
	UserID string  `json:"userId" binding:"required"`
	Type   int     `json:"type" binding:"required"`
	Amount int     `json:"amount" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

type OrderOutputDTO struct {
	ID     string  `json:"ID"`
	UserID string  `json:"userId"`
	Type   int     `json:"type"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
}

type OrderUsecaseInterface interface {
	ProcessOrder(orderInput *OrderInputDTO) error
}

type orderUsecase struct {
	orderRepositoryInterface  order.OrderRepositoryInterface
	queueInterface            order.OrderQueueInterface
	walletRepositoryInterface wallet.WalletRepositoryInterface
}

func NewOrderUsecase(
	orderRepositoryInterface order.OrderRepositoryInterface,
	walletRepositoryInterface wallet.WalletRepositoryInterface,
	queueInterface order.OrderQueueInterface,
) *orderUsecase {
	orderUsecaseObj := &orderUsecase{
		orderRepositoryInterface,
		queueInterface,
		walletRepositoryInterface}
	orderUsecaseObj.StartOrderProcessingWorker()
	return orderUsecaseObj
}

func (ou *orderUsecase) StartOrderProcessingWorker() {
	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		numWorkers = 10
	}

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				order, err := ou.queueInterface.DequeueOrder()
				if err != nil {
					continue
				}

				err = ou.ExecuteOrder(order)
				if err != nil {
					log.Printf("Failed to process order: %v", err)
				}
			}
		}()
	}
}

func (ou *orderUsecase) ProcessOrder(orderInput *OrderInputDTO) error {
	orderEntity, err := order.NewOrder(orderInput.UserID, orderInput.Type, orderInput.Amount, orderInput.Price)
	if err != nil {
		return err
	}

	return ou.queueInterface.EnqueueOrder(orderEntity)
}

func (ou *orderUsecase) ExecuteOrder(orderEntity *order.Order) error {
	switch orderEntity.Type {
	case order.OrderTypeBuy:
		wallet, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
		if err != nil {
			return err
		}
		if wallet.Balance < float64(orderEntity.Amount)*orderEntity.Price {
			return errors.New("insufficient balance")
		}

		wallet.Balance -= float64(orderEntity.Amount) * orderEntity.Price
		wallet.Vibranium += orderEntity.Amount

		err = ou.walletRepositoryInterface.UpdateWallet(wallet)
		if err != nil {
			return err
		}

	case order.OrderTypeSell:
		wallet, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
		if err != nil {
			return err
		}

		if wallet.Vibranium < orderEntity.Amount {
			return errors.New("insufficient vibranium")
		}

		wallet.Vibranium -= orderEntity.Amount
		wallet.Balance += float64(orderEntity.Amount) * orderEntity.Price
		err = ou.walletRepositoryInterface.UpdateWallet(wallet)
		if err != nil {
			return err
		}

	default:
		orderEntity.Status = order.OrderStatusCompleted
		err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
		if err != nil {
			return err
		}
	}

	return nil
}
