package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

type WalletRepository struct {
	wallets map[string]*wallet.Wallet
}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{
		wallets: make(map[string]*wallet.Wallet),
	}
}

func (wr *WalletRepository) CreateWallet(wallet *wallet.Wallet) error {
	wr.UpdateLocalWalletReference(wallet)

	return nil
}
