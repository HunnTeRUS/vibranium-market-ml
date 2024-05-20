package wallet_controller

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WalletController struct {
	walletUsecase wallet_usecase.WalletUsecaseInterface
}

func NewWalletController(
	walletUsecase wallet_usecase.WalletUsecaseInterface) *WalletController {
	return &WalletController{walletUsecase}
}

func (w *WalletController) CreateWalletController(c *gin.Context) {
	wallet, err := w.walletUsecase.CreateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}
