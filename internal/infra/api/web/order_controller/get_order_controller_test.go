package order_controller

import (
	"errors"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderController_GetOrder(t *testing.T) {
	t.Run("GetOrder - successful order retrieval", func(t *testing.T) {
		mockUsecase := new(MockOrderUsecase)
		controller := NewOrderController(mockUsecase)

		orderOutput := &order_usecase.OrderOutputDTO{
			ID:     "order1",
			UserID: "user1",
			Type:   1,
			Amount: 100,
			Price:  50.0,
			Status: "completed",
		}
		mockUsecase.On("GetOrder", "order1").Return(orderOutput, nil)

		router := gin.Default()
		router.GET("/order/:id", controller.GetOrder)

		req, _ := http.NewRequest(http.MethodGet, "/order/order1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("GetOrder - order not found", func(t *testing.T) {
		mockUsecase := new(MockOrderUsecase)
		controller := NewOrderController(mockUsecase)

		mockUsecase.On("GetOrder", "order1").Return(nil, errors.New(
			fmt.Sprintf("order %s not found", "order1")))

		router := gin.Default()
		router.GET("/order/:id", controller.GetOrder)

		req, _ := http.NewRequest(http.MethodGet, "/order/order1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}
