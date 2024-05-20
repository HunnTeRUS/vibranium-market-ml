package order_controller

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/usecase/order_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderController struct {
	orderUsecase order_usecase.OrderUsecaseInterface
}

func NewOrderController(orderUsecase order_usecase.OrderUsecaseInterface) *OrderController {
	return &OrderController{orderUsecase}
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var orderInput order_usecase.OrderInputDTO

	if err := c.ShouldBindJSON(&orderInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := oc.orderUsecase.ProcessOrder(&orderInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
