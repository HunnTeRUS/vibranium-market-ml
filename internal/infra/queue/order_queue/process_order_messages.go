package order_queue

import (
	"encoding/json"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

type orderQueue struct {
	sqsConnection *sqs.SQS
	sqsQueueURL   string
}

func NewOrderQueue(sqsConnection *sqs.SQS) *orderQueue {
	return &orderQueue{
		sqsConnection,
		os.Getenv("SQS_QUEUE_URL")}
}

func (q *orderQueue) EnqueueOrder(order *order.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = q.sqsConnection.SendMessage(&sqs.SendMessageInput{
		QueueUrl:       aws.String(q.sqsQueueURL),
		MessageBody:    aws.String(string(orderJSON)),
		MessageGroupId: aws.String("default"),
	})
	return err
}

func (q *orderQueue) DequeueOrder() (*order.Order, error) {
	result, err := q.sqsConnection.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.sqsQueueURL),
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return nil, err
	}

	if len(result.Messages) == 0 {
		return nil, nil
	}

	var order order.Order
	err = json.Unmarshal([]byte(*result.Messages[0].Body), &order)
	if err != nil {
		return nil, err
	}

	_, err = q.sqsConnection.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.sqsQueueURL),
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})
	if err != nil {
		return nil, err
	}

	return &order, nil
}
