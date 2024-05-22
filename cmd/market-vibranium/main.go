package main

import (
	"database/sql"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/database/mysql"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/order_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/wallet_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/queue/order_queue"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/order_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/wallet_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.Info("About to start application")
	godotenv.Load("cmd/market-vibranium/.env")

	gin.SetMode(gin.ReleaseMode)

	dbClient := mysql.InitDB()
	defer dbClient.Close()

	walletController, orderController, orderQueue := initDependencies(dbClient)

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
			logger.Error("Failed to save snapshot:", err)
		} else {
			logger.Info("Snapshot saved successfully.")
		}

		os.Exit(0)
	}()

	logger.Info("Application up and running")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func initDependencies(
	dbConnection *sql.DB) (
	walletController *wallet_controller.WalletController,
	orderController *order_controller.OrderController,
	orderQueue *order_queue.OrderQueue) {

	walletRepository := wallet_repository.NewWalletRepository(dbConnection)
	orderQueue = order_queue.NewOrderQueue(1024)

	orderUsecase := order_usecase.NewOrderUsecase(
		order_repository.NewOrderRepository(dbConnection),
		walletRepository,
		orderQueue)

	walletController = wallet_controller.NewWalletController(
		wallet_usecase.NewWalletUsecase(walletRepository))

	orderController = order_controller.NewOrderController(orderUsecase)

	return
}
