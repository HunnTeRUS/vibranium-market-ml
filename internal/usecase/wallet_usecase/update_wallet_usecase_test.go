package wallet_usecase

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateWallet(t *testing.T) {

	t.Run("got success when updating wallet", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)
		walletEntity := &wallet.Wallet{
			UserID:    uuid.New().String(),
			Balance:   200.0,
			Vibranium: 20,
		}

		mockRepo.On("UpdateWallet", walletEntity).Return(nil)

		wu := NewWalletUsecase(mockRepo)

		err := wu.UpdateWallet(walletEntity)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error trying to update wallet", func(t *testing.T) {

		mockRepo := new(MockWalletRepository)
		walletEntity := &wallet.Wallet{
			UserID:    uuid.New().String(),
			Balance:   200.0,
			Vibranium: 20,
		}

		mockRepo.On("UpdateWallet", walletEntity).Return(errors.New("update error"))

		wu := NewWalletUsecase(mockRepo)

		err := wu.UpdateWallet(walletEntity)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
