package wallet

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWallet(t *testing.T) {
	t.Run("Create new wallet", func(t *testing.T) {
		wallet := NewWallet()
		assert.NotNil(t, wallet)
		assert.NotEmpty(t, wallet.UserID)
		assert.Equal(t, 0.0, wallet.Balance)
		assert.Equal(t, 0, wallet.Vibranium)
	})
}

func TestValidateDeposit(t *testing.T) {
	t.Run("Valid deposit", func(t *testing.T) {
		err := ValidateDeposit(uuid.New().String(), 100.0, 10)
		assert.NoError(t, err)
	})

	t.Run("Invalid userID", func(t *testing.T) {
		err := ValidateDeposit("invalid-uuid", 100.0, 10)
		assert.Error(t, err)
		assert.Equal(t, "invalid userId value", err.Error())
	})

	t.Run("Negative balance", func(t *testing.T) {
		err := ValidateDeposit(uuid.New().String(), -100.0, 10)
		assert.Error(t, err)
		assert.Equal(t, "it's not allowed to deposit less than 0 in balance value", err.Error())
	})

	t.Run("Negative vibranium", func(t *testing.T) {
		err := ValidateDeposit(uuid.New().String(), 100.0, -10)
		assert.Error(t, err)
		assert.Equal(t, "it's not allowed to deposit less than 0 in vibranium amount ", err.Error())
	})
}

func TestWallet_DebitVibranium(t *testing.T) {
	t.Run("Debit vibranium", func(t *testing.T) {
		wallet := &Wallet{Vibranium: 10}
		wallet.DebitVibranium(5)
		assert.Equal(t, 5, wallet.Vibranium)
	})
}

func TestWallet_CreditVibranium(t *testing.T) {
	t.Run("Credit vibranium", func(t *testing.T) {
		wallet := &Wallet{Vibranium: 10}
		wallet.CreditVibranium(5)
		assert.Equal(t, 15, wallet.Vibranium)
	})
}

func TestWallet_DebitBalance(t *testing.T) {
	t.Run("Debit balance", func(t *testing.T) {
		wallet := &Wallet{Balance: 100.0}
		wallet.DebitBalance(50.0)
		assert.Equal(t, 50.0, wallet.Balance)
	})
}

func TestWallet_CreditBalance(t *testing.T) {
	t.Run("Credit balance", func(t *testing.T) {
		wallet := &Wallet{Balance: 100.0}
		wallet.CreditBalance(50.0)
		assert.Equal(t, 150.0, wallet.Balance)
	})
}
