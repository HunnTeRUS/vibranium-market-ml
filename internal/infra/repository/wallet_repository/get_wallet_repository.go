package wallet_repository

import (
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

func (wr *WalletRepository) GetWalletBalance(userID string) (*wallet.Wallet, bool) {
	walletLocal, exists := wr.wallets[userID]
	return walletLocal, exists
}

func (wr *WalletRepository) GetWallet(userID string) (*wallet.Wallet, error) {
	v, ok := wr.GetWalletBalance(userID)
	if !ok {
		return nil, errors.New(fmt.Sprintf("wallet %s not found", userID))
	}

	return v, nil
}
