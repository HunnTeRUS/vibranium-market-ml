package wallet_repository

import (
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

func (wr *walletRepository) GetWalletBalance(userID string) (*wallet.Wallet, bool) {
	wr.Lock()
	defer wr.Unlock()
	walletLocal, exists := wr.wallets[userID]
	return walletLocal, exists
}

func (wr *walletRepository) GetWallet(userID string) (*wallet.Wallet, error) {
	if walletMemory, exists := wr.GetWalletBalance(userID); exists {
		return walletMemory, nil
	}

	tableName := os.Getenv("DYNAMODB_WALLETS_TABLE")
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}
	result, err := wr.dynamodbConnection.GetItem(input)
	if err != nil {
		logger.Error("error trying to get object from", err)
		return nil, err
	}
	if result.Item == nil {
		logger.Warn(fmt.Sprintf("wallet %s not found", userID))
		return nil, errors.New(fmt.Sprintf("wallet %s not found", userID))
	}
	wallet := new(wallet.Wallet)
	err = dynamodbattribute.UnmarshalMap(result.Item, wallet)
	if err != nil {
		logger.Error("error trying to unmarshal object for dynamodb", err)
		return nil, err
	}

	wr.UpdateWalletBalance(wallet)

	return wallet, nil
}
