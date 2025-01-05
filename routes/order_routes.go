package routes

import (
	"crud-clean-architecture/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.RouterGroup, handler *handler.OrderHandler) {
	r.POST("/", handler.CreateOrder)
	r.GET("/", handler.GetAllOrders)
	r.GET("/:id", handler.GetOrderByID)
	r.DELETE("/:id", handler.DeleteOrder)
}
