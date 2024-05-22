package wallet_repository

import (
	"database/sql"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"sync"
)

type walletRepository struct {
	sync.Mutex

	dynamodbConnection *dynamodb.DynamoDB
	wallets            map[string]*wallet.Wallet
	dbConnection       *sql.DB
}

func NewWalletRepository(dynamodbConnection *dynamodb.DynamoDB) *walletRepository {
	return &walletRepository{
		dynamodbConnection: dynamodbConnection,
		wallets:            make(map[string]*wallet.Wallet),
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
