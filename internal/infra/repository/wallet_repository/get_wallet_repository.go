package wallet_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

func (wr *walletRepository) GetWalletBalance(userID string) (*wallet.Wallet, bool) {
	wr.RLock()
	defer wr.RUnlock()
	walletLocal, exists := wr.wallets[userID]
	return walletLocal, exists
}

func (wr *walletRepository) GetWallet(userID string) (*wallet.Wallet, error) {
	stmt, err := wr.dbConnection.Prepare("SELECT * FROM wallet WHERE userId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)

	var wallet wallet.Wallet
	err = row.Scan(&wallet.UserID, &wallet.Balance, &wallet.Vibranium)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(fmt.Sprintf("wallet %s not found", userID))
			return nil, errors.New(fmt.Sprintf("wallet %s not found", userID))
		}
		return nil, err
	}

	wr.UpdateLocalWalletReference(&wallet)

	return &wallet, nil
}
