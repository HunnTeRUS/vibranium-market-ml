package order_repository

import (
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

func (u *OrderRepository) GetOrder(orderID string) (*order.Order, error) {
	tableName := os.Getenv("DYNAMODB_ORDERS_TABLE")

	fmt.Println(orderID)
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
		logger.Error("Error trying to get item from dynamodb", err)
		return nil, err
	}
	if result.Item == nil {
		logger.Warn(fmt.Sprintf("order %s was not found", orderID))
		return nil, errors.New("order %s not found")
	}

	orderDbEntity := new(OrderDynamoDBEntity)
	err = dynamodbattribute.UnmarshalMap(result.Item, orderDbEntity)
	if err != nil {
		logger.Error("Error trying to unmarshal object from dynamodb", err)
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
