package main

import (
	dynamodbConn "github.com/HunnTeRUS/vibranium-market-ml/config/database/dynamodb"
	redisConn "github.com/HunnTeRUS/vibranium-market-ml/config/queue/redis"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/order_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/wallet_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/queue/order_queue"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/order_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/wallet_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	godotenv.Load("cmd/market-vibranium/.env")

	gin.SetMode(gin.ReleaseMode)

	dynamodbClient := dynamodbConn.InitDB()

	redisClient := redisConn.InitQueue()

	walletController, orderController := initDependencies(redisClient, dynamodbClient)

	r := gin.Default()

	r.POST("/orders", orderController.CreateOrder)
	r.GET("/wallets/:userId", walletController.GetWallet)
	r.PATCH("/wallets/deposit", walletController.DepositToWallet)
	r.POST("/wallets", walletController.CreateWalletController)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func initDependencies(
	redisConnection *redis.Client, dynamoDbConnection *dynamodb.DynamoDB) (
	walletController *wallet_controller.WalletController,
	orderController *order_controller.OrderController) {

	walletRepository := wallet_repository.NewWalletRepository(dynamoDbConnection)

	orderUsecase := order_usecase.NewOrderUsecase(
		order_repository.NewOrderRepository(dynamoDbConnection),
		walletRepository,
		order_queue.NewOrderQueue(redisConnection))

	walletController = wallet_controller.NewWalletController(
		wallet_usecase.NewWalletUsecase(walletRepository))

	orderController = order_controller.NewOrderController(orderUsecase)

	return
}
