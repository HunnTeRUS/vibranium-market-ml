package wallet_controller

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (w *WalletController) DepositToWallet(c *gin.Context) {
	var depositWallet wallet_usecase.WalletDepositInputDTO
	if err := c.ShouldBindJSON(&depositWallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := w.walletUsecase.DepositWallet(depositWallet.UserID, depositWallet.Amount, depositWallet.Vibranium)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
