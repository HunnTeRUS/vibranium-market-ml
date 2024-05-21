package order_usecase

import "github.com/HunnTeRUS/vibranium-market-ml/config/logger"

func (ou *OrderUsecase) GetOrder(orderID string) (*OrderOutputDTO, error) {
	orderEntity, err := ou.orderRepositoryInterface.GetOrder(orderID)
	if err != nil {
		logger.Error("action=GetOrderUseCase, message=error from repository", err)
		return nil, err
	}

	return &OrderOutputDTO{
		ID:     orderEntity.ID,
		UserID: orderEntity.UserID,
		Type:   orderEntity.Type,
		Amount: orderEntity.Amount,
		Price:  orderEntity.Price,
		Status: orderEntity.Status,
	}, nil
}
