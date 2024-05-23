package wallet_controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalletController_DepositToWallet(t *testing.T) {
	t.Run("DepositToWallet - successful deposit", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		depositInput := wallet_usecase.WalletDepositInputDTO{
			UserID:    "user1",
			Amount:    50,
			Vibranium: 5,
		}
		mockUsecase.On("DepositWallet", depositInput.UserID, depositInput.Amount, depositInput.Vibranium).Return(nil)

		router := gin.Default()
		router.POST("/wallet/deposit", controller.DepositToWallet)

		body, _ := json.Marshal(depositInput)
		req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("DepositToWallet - invalid input", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		router := gin.Default()
		router.POST("/wallet/deposit", controller.DepositToWallet)

		invalidBody := []byte(`{invalid json}`)
		req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("DepositToWallet - deposit failure", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		depositInput := wallet_usecase.WalletDepositInputDTO{
			UserID:    "user1",
			Amount:    50,
			Vibranium: 5,
		}
		mockUsecase.On("DepositWallet", depositInput.UserID, depositInput.Amount, depositInput.Vibranium).Return(errors.New("deposit failed"))

		router := gin.Default()
		router.POST("/wallet/deposit", controller.DepositToWallet)

		body, _ := json.Marshal(depositInput)
		req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}
