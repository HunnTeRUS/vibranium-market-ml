package wallet_controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/wallet_usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletController_CreateWalletController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CreateWalletController - successful wallet creation", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		mockWallet := &wallet_usecase.WalletOuputDTO{UserID: "user1", Balance: 100, Vibranium: 10}
		mockUsecase.On("CreateWallet").Return(mockWallet, nil)

		router := gin.Default()
		router.POST("/wallet", controller.CreateWalletController)

		req, _ := http.NewRequest(http.MethodPost, "/wallet", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("CreateWalletController - wallet creation failure", func(t *testing.T) {
		mockUsecase := new(MockWalletUsecase)
		controller := NewWalletController(mockUsecase)

		mockUsecase.On("CreateWallet").Return(nil, errors.New("creation failed"))

		router := gin.Default()
		router.POST("/wallet", controller.CreateWalletController)

		req, _ := http.NewRequest(http.MethodPost, "/wallet", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}

type MockWalletUsecase struct {
	mock.Mock
}

func (m *MockWalletUsecase) CreateWallet() (*wallet_usecase.WalletOuputDTO, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*wallet_usecase.WalletOuputDTO), args.Error(1)
}

func (m *MockWalletUsecase) DepositWallet(userID string, amount float64, vibranium int) error {
	args := m.Called(userID, amount, vibranium)
	return args.Error(0)
}

func (m *MockWalletUsecase) GetWallet(userID string) (*wallet_usecase.WalletOuputDTO, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*wallet_usecase.WalletOuputDTO), args.Error(1)
}
