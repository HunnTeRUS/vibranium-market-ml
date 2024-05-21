package order_usecase

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/metrics"
	"os"
	"strconv"
)

type OrderInputDTO struct {
	UserID string  `json:"userId"`
	Type   int     `json:"type"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
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
				if err != nil || len(orderUnprocessed) == 0 {
					continue
				}

				metrics.ConcurrentOrdersProcessing.Inc()

				for _, orderV := range orderUnprocessed {
					if orderV == nil {
						continue
					}

					go func(orderV *order.Order) {
						defer metrics.ConcurrentOrdersProcessing.Dec()

						err = ou.ExecuteOrder(orderV)
						if err != nil {
							logger.Error("action=GetOrderUseCase, message=error trying to process order", err)
						}
					}(orderV)
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

	metrics.OrderPending.Inc()
	err = ou.orderRepositoryInterface.UpsertOrder(orderEntity)
	if err != nil {
		return "", err
	}

	if err := ou.queueInterface.EnqueueOrder(orderEntity); err != nil {
		return "", err
	}

	metrics.TotalValueProcessed.Add(orderInput.Price)

	return orderEntity.ID, nil
}

func (ou *OrderUsecase) ExecuteOrder(orderEntity *order.Order) error {
	defer metrics.OrderPending.Dec()

	switch orderEntity.Type {
	case order.OrderTypeBuy:
		walletEntity, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
		if err != nil {
			return err
		}
		if walletEntity.Balance < float64(orderEntity.Amount)*orderEntity.Price {
			orderEntity.Status = order.OrderStatusCanceled
			metrics.OrderCanceled.Inc()
			err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
			if err != nil {
				logger.Error("action=ExecuteOrder, message=error calling UpsertOrder repository for cancelling order", err)
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
			metrics.OrderCanceled.Inc()
			err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
			if err != nil {
				logger.Error("action=ExecuteOrder, message=error calling UpsertOrder repository for cancelling order", err)
				return err
			}

			return errors.New("insufficient vibranium")
		}

		walletEntity.Vibranium -= orderEntity.Amount
		walletEntity.Balance += float64(orderEntity.Amount) * orderEntity.Price
		err = ou.walletRepositoryInterface.UpdateWallet(walletEntity)
		if err != nil {
			logger.Error("action=ExecuteOrder, message=error calling UpdateWallet repository", err)
			return err
		}
	}

	orderEntity.Status = order.OrderStatusCompleted
	metrics.OrderProcessed.Inc()
	err := ou.orderRepositoryInterface.UpsertOrder(orderEntity)
	if err != nil {
		metrics.ProcessingErrors.Inc()
		logger.Error("action=ExecuteOrder, message=error calling UpsertOrder repository for completing order", err)
		return err
	}

	return nil
}
