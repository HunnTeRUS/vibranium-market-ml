package wallet_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

func (wr *walletRepository) UpdateWallet(wallet *wallet.Wallet) error {
	tableName := os.Getenv("DYNAMODB_WALLETS_TABLE")
	item, err := dynamodbattribute.MarshalMap(wallet)
	if err != nil {
		logger.Error("error trying to unmarshal object for dynamodb", err)
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}
	_, err = wr.dynamodbConnection.PutItem(input)
	return err
}
