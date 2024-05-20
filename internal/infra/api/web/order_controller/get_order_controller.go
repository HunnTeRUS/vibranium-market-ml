package order_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (oc *OrderController) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order, err := oc.orderUsecase.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
