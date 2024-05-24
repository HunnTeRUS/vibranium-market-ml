package wallet_usecase

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWallet(t *testing.T) {
	t.Run("get wallet returns with success", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)
		userID := uuid.New().String() // Gerando um UUID válido
		walletEntity := &wallet.Wallet{
			UserID:    userID,
			Balance:   100.0,
			Vibranium: 10,
		}

		mockRepo.On("GetWallet", userID).Return(walletEntity, nil)

		wu := NewWalletUsecase(mockRepo)

		result, err := wu.GetWallet(userID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, walletEntity.UserID, result.UserID)
		assert.Equal(t, walletEntity.Balance, result.Balance)
		assert.Equal(t, walletEntity.Vibranium, result.Vibranium)

		mockRepo.AssertExpectations(t)
	})

	t.Run("get wallet returns error", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)
		userID := uuid.New().String() // Gerando um UUID válido

		mockRepo.On("GetWallet", userID).Return((*wallet.Wallet)(nil), errors.New("get wallet error"))

		wu := NewWalletUsecase(mockRepo)

		result, err := wu.GetWallet(userID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "get wallet error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
