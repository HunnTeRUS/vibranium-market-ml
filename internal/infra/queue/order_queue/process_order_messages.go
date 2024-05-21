package order_queue

import (
	"encoding/json"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/order"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"time"
)

type orderQueue struct {
	sqsConnection *sqs.SQS
	sqsQueueURL   string
}

func NewOrderQueue(sqsConnection *sqs.SQS) *orderQueue {
	q := &orderQueue{
		sqsConnection: sqsConnection,
		sqsQueueURL:   os.Getenv("SQS_QUEUE_URL"),
	}

	const maxRetries = 100
	retries := 0

	for {
		_, err := q.sqsConnection.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			QueueUrl: aws.String(q.sqsQueueURL),
			AttributeNames: aws.StringSlice([]string{
				"All",
			}),
		})

		if err == nil {
			break
		}

		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			logger.Warn("Queue does not exist, waiting for it to be created...")
			time.Sleep(2 * time.Second)
			retries++
			if retries >= maxRetries {
				log.Fatalf("Queue was not created after %d retries, exiting.", maxRetries)
			}
			continue
		}

		log.Fatalf("Failed to check queue existence: %v", err)
	}

	logger.Info("Queue exists, proceeding...")
	return q
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
		logger.Error("Error trying to unmarshal message from queue", err)
		return nil, err
	}

	_, err = q.sqsConnection.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.sqsQueueURL),
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})
	if err != nil {
		logger.Error("Error trying to delete message from queue", err)
		return nil, err
	}

	return &order, nil
}
