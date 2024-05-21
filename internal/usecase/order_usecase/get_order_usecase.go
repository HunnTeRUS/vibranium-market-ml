package order_usecase

func (ou *OrderUsecase) GetOrder(orderID string) (*OrderOutputDTO, error) {
	orderEntity, err := ou.orderRepositoryInterface.GetOrder(orderID)
	if err != nil {
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
