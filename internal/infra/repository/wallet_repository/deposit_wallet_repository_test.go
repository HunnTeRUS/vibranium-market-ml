package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDepositToWallet(t *testing.T) {
	repo := NewWalletRepository()

	t.Run("wallet not found", func(t *testing.T) {
		err := repo.DepositToWallet("nonexistent", 100.0, 50)
		assert.Error(t, err)
		assert.Equal(t, "wallet nonexistent not found", err.Error())
	})

	t.Run("wallet found and deposit successful", func(t *testing.T) {
		expectedWallet := &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
		repo.wallets["user1"] = expectedWallet

		err := repo.DepositToWallet("user1", 50.0, 25)
		assert.NoError(t, err)

		repo.RLock()
		updatedWallet, exists := repo.wallets["user1"]
		repo.RUnlock()

		assert.True(t, exists)
		assert.Equal(t, 150.0, updatedWallet.Balance)
		assert.Equal(t, 75, updatedWallet.Vibranium)
	})
}
