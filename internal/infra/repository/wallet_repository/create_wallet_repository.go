package wallet_repository

import (
	"database/sql"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"sync"
)

type walletRepository struct {
	sync.RWMutex

	wallets      map[string]*wallet.Wallet
	dbConnection *sql.DB
}

func NewWalletRepository(dbConnection *sql.DB) *walletRepository {
	return &walletRepository{
		dbConnection: dbConnection,
		wallets:      make(map[string]*wallet.Wallet),
	}
}

func (wr *walletRepository) CreateWallet(wallet *wallet.Wallet) error {
	wr.UpdateLocalWalletReference(wallet)

	go func() {
		stmt, err := wr.dbConnection.Prepare("INSERT INTO wallet (userId, balance, vibranium) VALUES (?, ?, ?)")
		if err != nil {
			logger.Error("error trying to prepare database query", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(wallet.UserID, wallet.Balance, wallet.Vibranium)
		if err != nil {
			logger.Error("error trying to execute database query", err)
			return
		}
	}()

	return nil
}
