package wallet_repository

import (
	"errors"
	"fmt"
)

func (wr *WalletRepository) DepositToWallet(userID string, amount float64, vibranium int) error {
	if wallet, exists := wr.GetWalletBalance(userID); exists {
		wallet.Balance += amount
		wallet.Vibranium += vibranium

		return wr.UpdateWallet(wallet)
	}

	return errors.New(fmt.Sprintf("wallet %s not found", userID))
}
