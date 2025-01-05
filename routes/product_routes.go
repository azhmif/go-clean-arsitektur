package routes

import (
	"crud-clean-architecture/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.RouterGroup, handler *handler.ProductHandler) {
	r.POST("/", handler.CreateProduct)
	r.GET("/", handler.GetAllProducts)
	r.GET("/:id", handler.GetProductByID)
	r.PUT("/:id", handler.UpdateProduct)
	r.DELETE("/:id", handler.DeleteProduct)
}
