package main

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	godotenv.Load("cmd/market-vibranium/.env")

	region := os.Getenv("AWS_REGION")
	walletsTable := os.Getenv("DYNAMODB_WALLETS_TABLE")
	ordersTable := os.Getenv("DYNAMODB_ORDERS_TABLE")
	dynamodbEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(dynamodbEndpoint),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	}))

	svc := dynamodb.New(sess)

	// Create Wallets table
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(walletsTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("UserID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("UserID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	// Create Orders table
	_, err = svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(ordersTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	log.Println("Tables created successfully")
}
