package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWalletRepository(t *testing.T) {
	repo := NewWalletRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.wallets)
	assert.Equal(t, 0, len(repo.wallets))
}

func TestCreateWallet(t *testing.T) {
	repo := NewWalletRepository()

	t.Run("create new wallet", func(t *testing.T) {
		w := &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
		err := repo.CreateWallet(w)
		assert.NoError(t, err)

		repo.RLock()
		createdWallet, exists := repo.wallets["user1"]
		repo.RUnlock()

		assert.True(t, exists)
		assert.Equal(t, w, createdWallet)
	})

	t.Run("create another wallet", func(t *testing.T) {
		w := &wallet.Wallet{UserID: "user2", Balance: 200.0, Vibranium: 100}
		err := repo.CreateWallet(w)
		assert.NoError(t, err)

		repo.RLock()
		createdWallet, exists := repo.wallets["user2"]
		repo.RUnlock()

		assert.True(t, exists)
		assert.Equal(t, w, createdWallet)
	})
}
