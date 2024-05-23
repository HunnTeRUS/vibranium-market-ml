package order_controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockOrderUsecase struct {
	mock.Mock
}

func (m *MockOrderUsecase) ProcessOrder(orderInput *order_usecase.OrderInputDTO) (string, error) {
	args := m.Called(orderInput)
	return args.String(0), args.Error(1)
}

func (m *MockOrderUsecase) GetOrder(orderID string) (*order_usecase.OrderOutputDTO, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order_usecase.OrderOutputDTO), args.Error(1)
}

func TestOrderController_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CreateOrder - successful order creation", func(t *testing.T) {
		mockUsecase := new(MockOrderUsecase)
		controller := NewOrderController(mockUsecase)

		orderInput := &order_usecase.OrderInputDTO{
			UserID: "user1",
			Type:   1,
			Amount: 100,
			Price:  50.0,
		}
		mockUsecase.On("ProcessOrder", orderInput).Return("order1", nil)

		router := gin.Default()
		router.POST("/order", controller.CreateOrder)

		body, _ := json.Marshal(orderInput)
		req, _ := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("CreateOrder - invalid input", func(t *testing.T) {
		mockUsecase := new(MockOrderUsecase)
		controller := NewOrderController(mockUsecase)

		router := gin.Default()
		router.POST("/order", controller.CreateOrder)

		invalidBody := []byte(`{invalid json}`)
		req, _ := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("CreateOrder - order creation failure", func(t *testing.T) {
		mockUsecase := new(MockOrderUsecase)
		controller := NewOrderController(mockUsecase)

		orderInput := &order_usecase.OrderInputDTO{
			UserID: "user1",
			Type:   1,
			Amount: 100,
			Price:  50.0,
		}
		mockUsecase.On("ProcessOrder", orderInput).Return("", errors.New("creation failed"))

		router := gin.Default()
		router.POST("/order", controller.CreateOrder)

		body, _ := json.Marshal(orderInput)
		req, _ := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}
