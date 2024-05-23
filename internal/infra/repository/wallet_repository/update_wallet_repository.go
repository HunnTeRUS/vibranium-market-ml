package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

func (wr *WalletRepository) UpdateWallet(wallet *wallet.Wallet) error {
	wr.UpdateLocalWalletReference(wallet)

	return nil
}

func (wr *WalletRepository) UpdateLocalWalletReference(wallet *wallet.Wallet) {
	wr.Lock()
	if wallet != nil {
		wr.wallets[wallet.UserID] = wallet
	}
	wr.Unlock()
}
