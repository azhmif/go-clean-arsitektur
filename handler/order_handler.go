package handler

import (
	"net/http"
	"strconv"

	"crud-clean-architecture/domain"
	"crud-clean-architecture/service"
	"crud-clean-architecture/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order domain.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid input", nil, validationErrors)
		return
	}

	err := h.orderService.CreateOrder(&order)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Order created successfully", order, nil)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "Failed to fetch orders", nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Orders fetched successfully", orders, nil)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Order fetched successfully", order, nil)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.orderService.DeleteOrder(uint(id))
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Order deleted successfully", nil, nil)
}
