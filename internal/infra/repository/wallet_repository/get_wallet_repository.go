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

func (wr *walletRepository) GetWallet(userId string) (*wallet.Wallet, error) {
	tableName := os.Getenv("DYNAMODB_WALLETS_TABLE")
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userId),
			},
		},
	}
	result, err := wr.dynamodbConnection.GetItem(input)
	if err != nil {
		logger.Error("error trying to get object from", err)
		return nil, err
	}
	if result.Item == nil {
		logger.Warn(fmt.Sprintf("wallet %s not found", userId))
		return nil, errors.New(fmt.Sprintf("wallet %s not found", userId))
	}
	wallet := new(wallet.Wallet)
	err = dynamodbattribute.UnmarshalMap(result.Item, wallet)
	if err != nil {
		logger.Error("error trying to unmarshal object for dynamodb", err)
		return nil, err
	}
	return wallet, nil
}
