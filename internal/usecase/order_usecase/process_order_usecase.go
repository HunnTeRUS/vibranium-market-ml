package order_usecase

import (
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
	Symbol string  `json:"symbol"`
}

type OrderOutputDTO struct {
	ID     string  `json:"orderId"`
	UserID string  `json:"userId"`
	Type   int     `json:"type"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
	Symbol string  `json:"symbol"`
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
				if err != nil || orderUnprocessed == nil {
					continue
				}

				metrics.ConcurrentOrdersProcessing.Inc()

				go func(orderV *order.Order) {
					defer metrics.ConcurrentOrdersProcessing.Dec()

					err = ou.ExecuteOrder(orderV)
					if err != nil {
						logger.Error("action=StartOrderProcessingWorker, message=error trying to process order", err)
					}
				}(orderUnprocessed)
			}
		}()
	}
}

func (ou *OrderUsecase) ProcessOrder(orderInput *OrderInputDTO) (string, error) {
	orderEntity, err := order.NewOrder(
		orderInput.UserID,
		orderInput.Type,
		orderInput.Amount,
		orderInput.Price,
		orderInput.Symbol)
	if err != nil {
		return "", err
	}

	metrics.OrderPending.Inc()
	ou.orderRepositoryInterface.UpsertOrder(orderEntity)

	if err := ou.queueInterface.EnqueueOrder(orderEntity); err != nil {
		return "", err
	}

	metrics.TotalValueProcessed.Add(orderInput.Price)

	return orderEntity.ID, nil
}

func (ou *OrderUsecase) ExecuteOrder(orderEntity *order.Order) error {
	defer metrics.OrderPending.Dec()

	matchOrder, err := ou.findMatchingOrder(orderEntity)
	if err != nil {
		return err
	}

	if matchOrder != nil {
		buyerWallet, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
		if err != nil {
			return err
		}

		sellerWallet, err := ou.walletRepositoryInterface.GetWallet(matchOrder.UserID)
		if err != nil {
			return err
		}

		if orderEntity.Type == order.OrderTypeBuy {
			if buyerWallet.Balance < float64(orderEntity.Amount)*orderEntity.Price {
				return orderEntity.CancelOrder(ou.orderRepositoryInterface, "insufficient balance")
			}

			buyerWallet.DebitBalance(float64(orderEntity.Amount) * orderEntity.Price)
			sellerWallet.CreditBalance(float64(orderEntity.Amount) * orderEntity.Price)

			buyerWallet.CreditVibranium(orderEntity.Amount)
			sellerWallet.DebitVibranium(orderEntity.Amount)
		} else if orderEntity.Type == order.OrderTypeSell {
			if sellerWallet.Vibranium < orderEntity.Amount {
				return orderEntity.CancelOrder(ou.orderRepositoryInterface, "insufficient vibranium")
			}

			sellerWallet.DebitVibranium(orderEntity.Amount)
			buyerWallet.CreditVibranium(orderEntity.Amount)

			sellerWallet.CreditBalance(float64(orderEntity.Amount) * orderEntity.Price)
			buyerWallet.DebitBalance(float64(orderEntity.Amount) * orderEntity.Price)
		}

		err = ou.walletRepositoryInterface.UpdateWallet(buyerWallet)
		if err != nil {
			logger.Error("action=ExecuteOrder, message=error calling UpdateWallet repository", err)
			return err
		}

		err = ou.walletRepositoryInterface.UpdateWallet(sellerWallet)
		if err != nil {
			logger.Error("action=ExecuteOrder, message=error calling UpdateWallet repository", err)
			return err
		}

		err = matchOrder.CompleteOrder(ou.orderRepositoryInterface)
		if err != nil {
			logger.Error("action=ExecuteOrder, message=error calling CompleteOrder repository", err)
			return err
		}

		err = orderEntity.CompleteOrder(ou.orderRepositoryInterface)
		if err != nil {
			logger.Error("action=ExecuteOrder, message=error calling CompleteOrder repository", err)
			return err
		}
	}

	return nil
}

func (ou *OrderUsecase) findMatchingOrder(orderEntity *order.Order) (*order.Order, error) {
	matchType := 0
	if orderEntity.Type == order.OrderTypeBuy {
		matchType = order.OrderTypeSell
	} else {
		matchType = order.OrderTypeBuy
	}

	orders, err := ou.orderRepositoryInterface.GetPendingOrders(orderEntity.Symbol, matchType)
	if err != nil {
		return nil, err
	}

	for _, o := range orders {
		if o.UserID == orderEntity.UserID {
			continue
		}

		if orderEntity.Price == o.Price && orderEntity.Amount == o.Amount {
			return o, nil
		}
	}

	return nil, nil
}
