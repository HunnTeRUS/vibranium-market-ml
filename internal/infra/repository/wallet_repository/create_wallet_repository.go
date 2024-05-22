package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"sync"
)

type walletRepository struct {
	sync.Mutex

	dynamodbConnection *dynamodb.DynamoDB
	wallets            map[string]*wallet.Wallet
}

func NewWalletRepository(dynamodbConnection *dynamodb.DynamoDB) *walletRepository {
	return &walletRepository{
		dynamodbConnection: dynamodbConnection,
		wallets:            make(map[string]*wallet.Wallet),
	}
}

func (wr *walletRepository) CreateWallet(wallet *wallet.Wallet) error {
	wr.UpdateWalletBalance(wallet)

	go func() {
		tableName := os.Getenv("DYNAMODB_WALLETS_TABLE")
		item, err := dynamodbattribute.MarshalMap(wallet)
		if err != nil {
			logger.Error("error trying to marshal object for dynamodb", err)
		}
		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		}
		_, err = wr.dynamodbConnection.PutItem(input)
		if err != nil {
			logger.Error("error trying to put item in dynamodb", err)
		}
	}()

	return nil
}
