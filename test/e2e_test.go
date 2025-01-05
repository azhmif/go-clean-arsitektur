package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"crud-clean-architecture/config"
	"crud-clean-architecture/domain"
	"crud-clean-architecture/handler"
	"crud-clean-architecture/repository"
	"crud-clean-architecture/routes"
	"crud-clean-architecture/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}
	// Setup database
	db := config.InitDB()
	_ = db.AutoMigrate(&domain.Category{}, &domain.Product{}, &domain.Order{}, &domain.OrderDetail{})

	// Setup Redis
	config.InitRedis()
	redisClient := config.RedisClient

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	productRepo := repository.NewProductRepository(db, redisClient)
	orderRepo := repository.NewOrderRepository(db, redisClient)

	// Initialize services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)

	// Initialize handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup router
	r := gin.Default()
	routes.RegisterCategoryRoutes(r.Group("/categories"), categoryHandler)
	routes.RegisterProductRoutes(r.Group("/products"), productHandler)
	routes.RegisterOrderRoutes(r.Group("/orders"), orderHandler)

	return r
}

func TestE2EOrderAPI(t *testing.T) {
	router := setupTestRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Step 1: Create a new category
	categoryPayload := map[string]string{"name": "Electronics"}
	categoryBody, _ := json.Marshal(categoryPayload)
	resp, err := http.Post(server.URL+"/categories", "application/json", bytes.NewBuffer(categoryBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Step 2: Create a new product
	productPayload := map[string]interface{}{
		"name":        "Laptop",
		"price":       1500.00,
		"category_id": 1,
	}
	productBody, _ := json.Marshal(productPayload)
	resp, err = http.Post(server.URL+"/products", "application/json", bytes.NewBuffer(productBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Step 3: Create a new order
	orderPayload := map[string]interface{}{
		"invoice": "INV-001",
		"details": []map[string]interface{}{
			{"product_id": 1, "quantity": 2},
		},
	}
	orderBody, _ := json.Marshal(orderPayload)
	resp, err = http.Post(server.URL+"/orders", "application/json", bytes.NewBuffer(orderBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Step 4: Get the order
	resp, err = http.Get(server.URL + "/orders/1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3000), data["total_price"])
}
