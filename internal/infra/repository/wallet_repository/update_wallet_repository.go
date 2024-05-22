package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

func (wr *walletRepository) UpdateWallet(wallet *wallet.Wallet) error {
	wr.UpdateLocalWalletReference(wallet)

	go func() {
		stmt, err := wr.dbConnection.Prepare("UPDATE wallet SET balance = ?, vibranium = ? WHERE userId = ?")
		if err != nil {
			logger.Error("error trying to prepare update statement", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(wallet.Balance, wallet.Vibranium, wallet.UserID)
		if err != nil {
			logger.Error("error trying to update wallet", err)
			return
		}
	}()

	return nil
}

func (wr *walletRepository) UpdateLocalWalletReference(wallet *wallet.Wallet) {
	wr.Lock()
	if wallet != nil {
		wr.wallets[wallet.UserID] = wallet
	}
	wr.Unlock()
}
