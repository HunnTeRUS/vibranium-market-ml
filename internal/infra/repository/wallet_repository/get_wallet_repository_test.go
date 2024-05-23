package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWalletBalance(t *testing.T) {
	repo := NewWalletRepository()

	w, exists := repo.GetWalletBalance("nonexistent")
	assert.Nil(t, w)
	assert.False(t, exists)

	expectedWallet := &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
	repo.wallets["user1"] = expectedWallet

	w, exists = repo.GetWalletBalance("user1")
	assert.NotNil(t, w)
	assert.True(t, exists)
	assert.Equal(t, expectedWallet, w)
}

func TestGetWallet(t *testing.T) {
	repo := NewWalletRepository()

	_, err := repo.GetWallet("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "wallet nonexistent not found", err.Error())

	expectedWallet := &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
	repo.wallets["user1"] = expectedWallet

	w, err := repo.GetWallet("user1")
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.Equal(t, expectedWallet, w)
}
