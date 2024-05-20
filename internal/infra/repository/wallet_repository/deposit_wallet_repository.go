package wallet_repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func (wr *walletRepository) DepositToWallet(userID string, amount float64) error {
	tableName := os.Getenv("DYNAMODB_WALLETS_TABLE")
	key := map[string]*dynamodb.AttributeValue{
		"UserID": {
			S: aws.String(userID),
		},
	}
	update := map[string]*dynamodb.AttributeValueUpdate{
		"Balance": {
			Action: aws.String("ADD"),
			Value: &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%f", amount)),
			},
		},
	}
	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(tableName),
		Key:              key,
		AttributeUpdates: update,
		ReturnValues:     aws.String("UPDATED_NEW"),
	}
	_, err := wr.dynamodbConnection.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
