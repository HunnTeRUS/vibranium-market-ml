package order_usecase

import "github.com/HunnTeRUS/vibranium-market-ml/config/logger"

func (ou *OrderUsecase) GetOrder(orderID string) (*OrderOutputDTO, error) {
	if value, exists := ou.orderRepositoryInterface.GetMemOrder(orderID); exists {
		return &OrderOutputDTO{
			ID:     value.ID,
			UserID: value.UserID,
			Type:   value.Type,
			Amount: value.Amount,
			Price:  value.Price,
			Status: value.Status,
		}, nil
	}

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
