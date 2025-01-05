package handler

import (
	"net/http"
	"strconv"

	"crud-clean-architecture/domain"
	"crud-clean-architecture/service"
	"crud-clean-architecture/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.ProductForm
	// Validasi input
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		utils.JSONResponse(c, http.StatusBadRequest, "Validation error", nil, validationErrors)
		return
	}

	// Validasi keunikan nama kategori
	isUnique, err := h.productService.IsProductNameUnique(req.Name, req.CategoryID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "Failed to check product uniqueness", nil, nil)
		return
	}
	if !isUnique {
		// Format validation error untuk keunikan
		validationErrors := map[string]string{
			"name": "name must be unique",
		}
		utils.JSONResponse(c, http.StatusBadRequest, "Validation error", nil, validationErrors)
		return
	}

	// Simpan produk baru
	product := domain.Product{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}
	if err := h.productService.CreateProduct(&product); err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Product created successfully", product, nil)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "Failed to fetch products", nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Products retrieved successfully", products, nil)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid ID format", nil, nil)
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, "Product not found", nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Product retrieved successfully", product, nil)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid ID format", nil, nil)
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid request payload", nil, nil)
		return
	}

	product.ID = uint(id)
	if err := h.productService.UpdateProduct(&product); err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "Failed to update product", nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Product updated successfully", product, nil)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid ID format", nil, nil)
		return
	}

	err = h.productService.DeleteProduct(uint(id))
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Product deleted successfully", nil, nil)
}
