package main

import (
	"log"
	"reflect"

	"crud-clean-architecture/config"
	"crud-clean-architecture/domain"
	"crud-clean-architecture/handler"
	"crud-clean-architecture/middleware"
	"crud-clean-architecture/repository"
	"crud-clean-architecture/routes"
	"crud-clean-architecture/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}
	config.InitRedis()
	redisClient := config.RedisClient
	// Setup Database
	db := config.InitDB()

	// Migrate Database
	err := db.AutoMigrate(&domain.Category{}, &domain.Product{}, &domain.Order{}, &domain.OrderDetail{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Initialize Repositories
	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	productRepo := repository.NewProductRepository(db, redisClient)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize Services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)

	// Initialize Handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Router
	r := gin.Default()

	// Setup custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("json")
		})
	}

	// Add middleware
	r.Use(middleware.ErrorMiddleware())

	// Register Routes
	routes.RegisterCategoryRoutes(r.Group("/categories"), categoryHandler)
	routes.RegisterProductRoutes(r.Group("/products"), productHandler)
	routes.RegisterOrderRoutes(r.Group("/orders"), orderHandler)

	// Run the Server
	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
