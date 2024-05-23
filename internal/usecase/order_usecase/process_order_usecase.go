package order_usecase

import (
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/metrics"
)

type OrderInputDTO struct {
	UserID string  `json:"userId"`
	Type   int     `json:"type"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
}

type OrderOutputDTO struct {
	ID              string  `json:"orderId"`
	UserID          string  `json:"userId"`
	Type            int     `json:"type"`
	Amount          int     `json:"amount"`
	Price           float64 `json:"price"`
	Status          string  `json:"status"`
	AmountRemaining int     `json:"amountRemaining,omitempty"`
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

	go orderUsecaseObj.StartOrderProcessingWorker()

	return orderUsecaseObj
}
func (ou *OrderUsecase) StartOrderProcessingWorker() {
	for {
		go func() {
			orderUnprocessed, err := ou.queueInterface.DequeueOrder()
			if err != nil || orderUnprocessed == nil {
				return
			}

			metrics.ConcurrentOrdersProcessing.Inc()

			defer metrics.ConcurrentOrdersProcessing.Dec()

			err = ou.ExecuteOrder(orderUnprocessed)
			if err != nil {
				logger.Error("action=StartOrderProcessingWorker, message=error trying to process order", err)
			}
		}()
	}
}

func (ou *OrderUsecase) ProcessOrder(orderInput *OrderInputDTO) (string, error) {
	orderEntity, err := order.NewOrder(
		orderInput.UserID,
		orderInput.Type,
		orderInput.Amount,
		orderInput.Price)
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

	actualWallet, err := ou.validateCurrentWalletAndOrder(orderEntity)
	if err != nil {
		return err
	}

	var matchOrder *order.Order
	if orderEntity.Type == order.OrderTypeBuy {
		matchOrder, err = ou.orderRepositoryInterface.GetSellingMatchingOrder(orderEntity)
		if err != nil || matchOrder == nil {
			return err
		}

		matchWallet, err := ou.walletRepositoryInterface.GetWallet(matchOrder.UserID)
		if err != nil {
			return err
		}

		if matchOrder.SellValueRemaining > orderEntity.Amount {
			matchOrder.SellValueRemaining -= orderEntity.Amount
			subOrder, _ := order.NewOrder(
				matchOrder.UserID, matchOrder.Type, orderEntity.Amount, matchOrder.Price)

			subOrder.CorrelationId = matchOrder.ID

			subOrder.CompleteOrder(ou.orderRepositoryInterface)
			ou.orderRepositoryInterface.UpsertOrder(matchOrder)
			orderEntity.CompleteOrder(ou.orderRepositoryInterface)

		} else if matchOrder.SellValueRemaining == orderEntity.Amount {
			matchOrder.SellValueRemaining = 0
			matchOrder.CompleteOrder(ou.orderRepositoryInterface)
		}

		matchWallet.CreditBalance(float64(orderEntity.Amount) * orderEntity.Price)
		actualWallet.DebitBalance(float64(orderEntity.Amount) * orderEntity.Price)
		actualWallet.CreditVibranium(orderEntity.Amount)
		matchWallet.DebitVibranium(orderEntity.Amount)

		orderEntity.CompleteOrder(ou.orderRepositoryInterface)
	} else {
		for orderEntity.SellValueRemaining > 0 {
			matchOrder, err = ou.orderRepositoryInterface.GetBuyingMatchingOrder(orderEntity)
			if err != nil || matchOrder == nil {
				return err
			}

			matchWallet, err := ou.walletRepositoryInterface.GetWallet(matchOrder.UserID)
			if err != nil {
				return err
			}

			if orderEntity.SellValueRemaining > matchOrder.Amount {
				orderEntity.SellValueRemaining -= matchOrder.Amount
				subOrder, _ := order.NewOrder(
					orderEntity.UserID, orderEntity.Type, matchOrder.Amount, orderEntity.Price)

				subOrder.CorrelationId = orderEntity.ID

				ou.orderRepositoryInterface.UpsertOrder(orderEntity)
				subOrder.CompleteOrder(ou.orderRepositoryInterface)

			} else if matchOrder.Amount == orderEntity.SellValueRemaining {
				orderEntity.SellValueRemaining = 0
				orderEntity.CompleteOrder(ou.orderRepositoryInterface)
			}

			actualWallet.CreditBalance(float64(matchOrder.Amount) * matchOrder.Price)
			matchWallet.DebitBalance(float64(matchOrder.Amount) * matchOrder.Price)
			matchWallet.CreditVibranium(matchOrder.Amount)
			actualWallet.DebitVibranium(matchOrder.Amount)

			matchOrder.CompleteOrder(ou.orderRepositoryInterface)
		}
	}

	return nil
}

func (ou *OrderUsecase) validateCurrentWalletAndOrder(orderEntity *order.Order) (*wallet.Wallet, error) {
	actualWallet, err := ou.walletRepositoryInterface.GetWallet(orderEntity.UserID)
	if err != nil {
		return nil, err
	}

	switch orderEntity.Type {
	case order.OrderTypeBuy:
		if actualWallet.Balance < float64(orderEntity.Amount)*orderEntity.Price {
			return nil, orderEntity.CancelOrder(ou.orderRepositoryInterface,
				"insufficient balance within your wallet for buying this amount of vibranium for that price")
		}
	case order.OrderTypeSell:
		if actualWallet.Vibranium < orderEntity.Amount {
			return nil, orderEntity.CancelOrder(ou.orderRepositoryInterface,
				"insufficient vibranium within your wallet for selling")
		}
	}

	return actualWallet, nil
}
