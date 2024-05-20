package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

type orderRepository struct {
	dynamodbConnection *dynamodb.DynamoDB
}

func NewOrderRepository(dynamodbConnection *dynamodb.DynamoDB) *orderRepository {
	return &orderRepository{dynamodbConnection}
}

func (u *orderRepository) UpsertOrder(order *order.Order) error {
	tableName := os.Getenv("DYNAMODB_ORDERS_TABLE")
	item, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}
	_, err = u.dynamodbConnection.PutItem(input)
	return err
}
