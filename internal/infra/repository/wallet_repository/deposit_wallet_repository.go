package wallet_repository

import (
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
)

func (wr *walletRepository) DepositToWallet(userID string, amount float64) error {
	wallet, exists := wr.GetWalletBalance(userID)
	if !exists {
		return fmt.Errorf("wallet %s not found", userID)
	}

	wallet.Balance += amount

	err := wr.UpdateWallet(wallet)
	if err != nil {
		return err
	}

	go func() {
		stmt, err := wr.dbConnection.Prepare("UPDATE wallet SET balance = ? WHERE userId = ?")
		if err != nil {
			logger.Error("error trying to prepare sql statement", err)
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(wallet.Balance, userID)
		if err != nil {
			logger.Error("error trying to execute sql statement", err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			logger.Error("error trying to validate rows affected by database", err)
			return
		}

		if rowsAffected == 0 {
			logger.Error(fmt.Sprintf("wallet %s not found", userID), err)
			return
		}
	}()

	return nil
}
