package wallet_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (w *WalletController) GetWallet(c *gin.Context) {
	userId := c.Param("userId")
	wallet, err := w.walletUsecase.GetWallet(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}
