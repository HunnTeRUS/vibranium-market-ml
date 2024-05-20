package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func InitDB() *dynamodb.DynamoDB {
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	var sess *session.Session
	var err error

	if endpoint != "" {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Endpoint:    aws.String(endpoint),
			Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
		})
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	}

	if err != nil {
		panic(err)
	}

	return dynamodb.New(sess)
}
