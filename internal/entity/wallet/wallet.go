package wallet

import (
	"errors"
	"github.com/google/uuid"
)

type Wallet struct {
	UserID    string
	Balance   float64
	Vibranium int
}

func NewWallet() *Wallet {
	return &Wallet{
		UserID:    uuid.New().String(),
		Balance:   0,
		Vibranium: 0,
	}
}

func ValidateDeposit(userId string, amount float64) error {
	if uuid.Validate(userId) != nil {
		return errors.New("invalid userId value")
	}

	if amount <= 0 {
		return errors.New("it's not allowed to deposit 0 in value")
	}

	return nil
}

func (wl *Wallet) DebitVibranium(amount int) {
	wl.Vibranium -= amount
}

func (wl *Wallet) CreditVibranium(amount int) {
	wl.Vibranium += amount
}

func (wl *Wallet) DebitBalance(amount float64) {
	wl.Balance -= amount
}

func (wl *Wallet) CreditBalance(amount float64) {
	wl.Balance += amount
}

type WalletRepositoryInterface interface {
	CreateWallet(wallet *Wallet) error
	DepositToWallet(userID string, amount float64) error
	GetWallet(userId string) (*Wallet, error)
	UpdateWallet(wallet *Wallet) error
}
