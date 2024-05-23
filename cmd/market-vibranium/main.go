package main

import (
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/order_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/api/web/wallet_controller"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/queue/order_queue"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/order_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/infra/repository/wallet_repository"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
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
	gin.SetMode(gin.ReleaseMode)

	logger.Info("About to start application")
	godotenv.Load("cmd/market-vibranium/.env")

	walletRepository := wallet_repository.NewWalletRepository()
	orderQueue := order_queue.NewOrderQueue(20000)
	orderRepository := order_repository.NewOrderRepository()

	orderUsecase := order_usecase.NewOrderUsecase(
		orderRepository,
		walletRepository,
		orderQueue)

	walletController := wallet_controller.NewWalletController(
		wallet_usecase.NewWalletUsecase(walletRepository))

	orderController := order_controller.NewOrderController(orderUsecase)

	_ = orderQueue.LoadSnapshot()
	_ = orderRepository.LoadSnapshot()
	_ = walletRepository.LoadSnapshot()

	r := gin.New()

	r.POST("/orders", orderController.CreateOrder)
	r.GET("/orders/:id", orderController.GetOrder)
	r.GET("/wallets/:userId", walletController.GetWallet)
	r.PATCH("/wallets/deposit", walletController.DepositToWallet)
	r.POST("/wallets", walletController.CreateWalletController)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	go gracefullyShutdown(orderQueue, orderRepository, walletRepository)

	logger.Info("Application up and running")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func gracefullyShutdown(
	orderQueue *order_queue.OrderQueue,
	orderRepository *order_repository.OrderRepository,
	walletRepository *wallet_repository.WalletRepository,
) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	defer os.Exit(0)

	sig := <-sigChannel
	logger.Warn(fmt.Sprintf("Received signal: %s. Initiating graceful shutdown...\n", sig))

	logger.Warn("Waiting for 10 seconds to finish processing orders...")
	time.Sleep(10 * time.Second)

	err := orderQueue.SaveSnapshot()
	if err != nil {
		logger.Error("Failed to save snapshot:", err)
	}

	err = orderRepository.SaveSnapshot()
	if err != nil {
		logger.Error("Failed to save snapshot:", err)
		return
	}

	err = walletRepository.SaveSnapshot()
	if err != nil {
		logger.Error("Failed to save snapshot:", err)
		return
	}

	logger.Info("Snapshot saved successfully.")
}
