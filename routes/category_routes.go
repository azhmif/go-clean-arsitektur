package routes

import (
	"crud-clean-architecture/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(r *gin.RouterGroup, handler *handler.CategoryHandler) {
	r.POST("/", handler.CreateCategory)
	r.GET("/", handler.GetAllCategories)
	r.GET("/:id", handler.GetCategoryByID)
	r.PUT("/:id", handler.UpdateCategory)
	r.DELETE("/:id", handler.DeleteCategory)
}
