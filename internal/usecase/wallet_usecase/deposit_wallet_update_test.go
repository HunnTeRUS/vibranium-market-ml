package wallet_usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDepositWallet(t *testing.T) {
	t.Run("successfully deposit values inside wallet", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)
		userID := uuid.New().String() // Gerando um UUID válido
		amount := 100.0
		vibranium := 10

		mockRepo.On("DepositToWallet", userID, amount, vibranium).Return(nil)

		wu := NewWalletUsecase(mockRepo)

		err := wu.DepositWallet(userID, amount, vibranium)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error trying to deposit values inside wallet", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)
		userID := uuid.New().String() // Gerando um UUID válido
		amount := 100.0
		vibranium := 10

		mockRepo.On("DepositToWallet", userID, amount, vibranium).Return(errors.New("deposit error"))

		wu := NewWalletUsecase(mockRepo)

		err := wu.DepositWallet(userID, amount, vibranium)
		assert.Error(t, err)
		assert.Equal(t, "deposit error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
