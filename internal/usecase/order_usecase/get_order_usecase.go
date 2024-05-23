package order_usecase

import (
	"errors"
	"fmt"
)

func (ou *OrderUsecase) GetOrder(orderID string) (*OrderOutputDTO, error) {
	if value, exists := ou.orderRepositoryInterface.GetMemOrder(orderID); exists {
		return &OrderOutputDTO{
			ID:              value.ID,
			UserID:          value.UserID,
			Type:            value.Type,
			Amount:          value.Amount,
			Price:           value.Price,
			Status:          value.Status,
			AmountRemaining: value.SellValueRemaining,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("order %s not found", orderID))
}
