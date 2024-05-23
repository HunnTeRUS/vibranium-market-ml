package wallet_controller

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalletController_GetWallet(t *testing.T) {
	t.Run("GetWallet - successful retrieval", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		mockWallet := &wallet_usecase.WalletOuputDTO{UserID: "user1", Balance: 100, Vibranium: 10}
		mockUsecase.On("GetWallet", "user1").Return(mockWallet, nil)

		router := gin.Default()
		router.GET("/wallet/:userId", controller.GetWallet)

		req, _ := http.NewRequest(http.MethodGet, "/wallet/user1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("GetWallet - retrieval failure", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		mockUsecase.On("GetWallet", "user1").Return(nil, errors.New("retrieval failed"))

		router := gin.Default()
		router.GET("/wallet/:userId", controller.GetWallet)

		req, _ := http.NewRequest(http.MethodGet, "/wallet/user1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}
