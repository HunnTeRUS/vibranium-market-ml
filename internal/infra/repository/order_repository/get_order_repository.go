package order_repository

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

func (u *orderRepository) GetOrder(orderID string) (*order.Order, error) {
	tableName := os.Getenv("DYNAMODB_ORDERS_TABLE")
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(orderID),
			},
		},
	}

	result, err := u.dynamodbConnection.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New("order not found")
	}

	orderDbEntity := new(OrderDynamoDBEntity)
	err = dynamodbattribute.UnmarshalMap(result.Item, orderDbEntity)
	if err != nil {
		return nil, err
	}

	return &order.Order{
		ID:     orderDbEntity.ID,
		UserID: orderDbEntity.UserID,
		Type:   orderDbEntity.Type,
		Amount: orderDbEntity.Amount,
		Price:  orderDbEntity.Price,
		Status: orderDbEntity.Status,
	}, nil
}
