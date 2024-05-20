package order_queue

import (
	"context"
	"encoding/json"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/go-redis/redis/v8"
)

type orderQueue struct {
	redisConnection *redis.Client
}

func NewOrderQueue(redisConnection *redis.Client) *orderQueue {
	return &orderQueue{redisConnection}
}

func (q *orderQueue) EnqueueOrder(order *order.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return q.redisConnection.LPush(context.Background(), "orderQueue", orderJSON).Err()
}

func (q *orderQueue) DequeueOrder() (*order.Order, error) {
	result, err := q.redisConnection.RPop(context.Background(), "orderQueue").Result()
	if err != nil {
		return nil, err
	}
	var order order.Order
	err = json.Unmarshal([]byte(result), &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
