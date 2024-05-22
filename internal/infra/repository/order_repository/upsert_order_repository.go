package order_repository

import (
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"sync"
)

type OrderDynamoDBEntity struct {
	ID     string  `json:"ID" dynamodbav:"ID"`
	UserID string  `json:"userId" dynamodbav:"UserId"`
	Type   int     `json:"type" dynamodbav:"Type"`
	Amount int     `json:"amount" dynamodbav:"Amount"`
	Price  float64 `json:"price" dynamodbav:"Price"`
	Status string  `json:"status" dynamodbav:"Status"`
}

type OrderRepository struct {
	sync.Mutex
	dynamodbConnection *dynamodb.DynamoDB

	orders map[string]*order.Order
}

func NewOrderRepository(dynamodbConnection *dynamodb.DynamoDB) *OrderRepository {
	return &OrderRepository{dynamodbConnection: dynamodbConnection, orders: make(map[string]*order.Order)}
}

func (u *OrderRepository) UpsertOrder(order *order.Order) error {
	u.UpsertLocalOrder(order)

	go func() {
		tableName := os.Getenv("DYNAMODB_ORDERS_TABLE")
		item, err := dynamodbattribute.MarshalMap(&OrderDynamoDBEntity{
			ID:     order.ID,
			UserID: order.UserID,
			Type:   order.Type,
			Amount: order.Amount,
			Price:  order.Price,
			Status: order.Status,
		})

		if err != nil {
			logger.Error("Error trying to unmarshal object for dynamodb", err)
		}
		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		}
		_, err = u.dynamodbConnection.PutItem(input)
		if err != nil {
			logger.Error("Error trying to put object on dynamodb", err)
		}
	}()

	return nil
}

func (u *OrderRepository) UpsertLocalOrder(orderEntity *order.Order) {
	u.Lock()
	defer u.Unlock()
	u.orders[orderEntity.ID] = orderEntity
}
