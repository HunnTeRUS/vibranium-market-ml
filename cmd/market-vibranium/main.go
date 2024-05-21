package main

import (
	"fmt"
	dynamodbConn "github.com/HunnTeRUS/vibranium-market-ml/config/database/dynamodb"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/order_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/wallet_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/queue/order_queue"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/order_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/wallet_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	godotenv.Load("cmd/market-vibranium/.env")

	gin.SetMode(gin.ReleaseMode)

	dynamodbClient := dynamodbConn.InitDB()

	walletController, orderController, orderQueue := initDependencies(dynamodbClient)

	r := gin.Default()

	r.POST("/orders", orderController.CreateOrder)
	r.GET("/orders/:id", orderController.GetOrder)
	r.GET("/wallets/:userId", walletController.GetWallet)
	r.PATCH("/wallets/deposit", walletController.DepositToWallet)
	r.POST("/wallets", walletController.CreateWalletController)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChannel
		logger.Warn(fmt.Sprintf("Received signal: %s. Initiating graceful shutdown...\n", sig))

		logger.Warn("Waiting for 10 seconds to finish processing orders...")
		time.Sleep(10 * time.Second)

		err := orderQueue.SaveSnapshot()
		if err != nil {
			logger.Error("Failed to save snapshot to S3:", err)
		} else {
			logger.Info("Snapshot saved successfully.")
		}

		os.Exit(0)
	}()

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func initDependencies(
	dynamoDbConnection *dynamodb.DynamoDB) (
	walletController *wallet_controller.WalletController,
	orderController *order_controller.OrderController,
	orderQueue *order_queue.OrderQueue) {

	walletRepository := wallet_repository.NewWalletRepository(dynamoDbConnection)
	orderQueue = order_queue.NewOrderQueue(1024)

	orderUsecase := order_usecase.NewOrderUsecase(
		order_repository.NewOrderRepository(dynamoDbConnection),
		walletRepository,
		orderQueue)

	walletController = wallet_controller.NewWalletController(
		wallet_usecase.NewWalletUsecase(walletRepository))

	orderController = order_controller.NewOrderController(orderUsecase)

	return
}
