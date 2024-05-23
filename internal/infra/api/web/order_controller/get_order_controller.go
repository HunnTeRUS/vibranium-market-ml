package order_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (oc *OrderController) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order, err := oc.orderUsecase.GetOrder(orderID)
	if err != nil {
		if err.Error() == fmt.Sprintf("order %s not found", orderID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
