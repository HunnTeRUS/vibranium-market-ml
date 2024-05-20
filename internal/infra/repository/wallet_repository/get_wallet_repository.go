package wallet_repository

import (
	"errors"
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
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New("wallet not found")
	}
	wallet := new(wallet.Wallet)
	err = dynamodbattribute.UnmarshalMap(result.Item, wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
