package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUpdateWallet tests the UpdateWallet function
func TestUpdateWallet(t *testing.T) {
	repo := NewWalletRepository()

	w := &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
	err := repo.UpdateWallet(w)
	assert.NoError(t, err)

	updatedWallet, exists := repo.wallets[w.UserID]

	assert.True(t, exists)
	assert.Equal(t, w, updatedWallet)
}

// TestUpdateLocalWalletReference tests the UpdateLocalWalletReference function
func TestUpdateLocalWalletReference(t *testing.T) {
	repo := NewWalletRepository()

	w := &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
	repo.UpdateLocalWalletReference(w)

	updatedWallet, exists := repo.wallets[w.UserID]

	assert.True(t, exists)
	assert.Equal(t, w, updatedWallet)
}

// TestUpdateLocalWalletReference_NilWallet tests the UpdateLocalWalletReference function with a nil wallet
func TestUpdateLocalWalletReference_NilWallet(t *testing.T) {
	repo := NewWalletRepository()

	repo.UpdateLocalWalletReference(nil)

	// Ensure no panic and no update to the wallets map
	assert.Equal(t, 0, len(repo.wallets))
}
