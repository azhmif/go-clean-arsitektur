package handler

import (
	"net/http"
	"strconv"

	"crud-clean-architecture/domain"
	"crud-clean-architecture/service"
	"crud-clean-architecture/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CategoryForm

	// Validasi input
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		utils.JSONResponse(c, http.StatusBadRequest, "Validation error", nil, validationErrors)
		return
	}

	// Validasi keunikan nama kategori
	isUnique, err := h.categoryService.IsCategoryNameUnique(req.Name)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "Failed to check category uniqueness", nil, nil)
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

	// Simpan kategori baru
	category := domain.Category{
		Name: req.Name,
	}
	if err := h.categoryService.CreateCategory(&category); err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Category created successfully", category, nil)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Categories retrieved successfully", categories, nil)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid ID format", nil, nil)
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, "Category not found", nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Category retrieved successfully", category, nil)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid ID format", nil, nil)
		return
	}
	var req domain.CategoryForm

	// Validasi input
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		utils.JSONResponse(c, http.StatusBadRequest, "Validation error", nil, validationErrors)
		return
	}
	category := domain.Category{
		Name: req.Name,
	}
	category.ID = uint(id)
	if err := h.categoryService.UpdateCategory(&category); err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Category updated successfully", category, nil)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "Invalid ID format", nil, nil)
		return
	}
	// Hapus kategori
	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Category deleted successfully", nil, nil)
}
