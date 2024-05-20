package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

type OrderDynamoDBEntity struct {
	ID     string  `json:"ID"`
	UserID string  `json:"userId"`
	Type   int     `json:"type"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type orderRepository struct {
	dynamodbConnection *dynamodb.DynamoDB
}

func NewOrderRepository(dynamodbConnection *dynamodb.DynamoDB) *orderRepository {
	return &orderRepository{dynamodbConnection}
}

func (u *orderRepository) UpsertOrder(order *order.Order) error {
	tableName := os.Getenv("DYNAMODB_ORDERS_TABLE")
	item, err := dynamodbattribute.MarshalMap(&OrderDynamoDBEntity{
		ID:     order.ID,
		UserID: order.UserID,
		Type:   order.Type,
		Amount: order.Amount,
		Price:  order.Price,
	})
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
