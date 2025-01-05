package service

import (
	"crud-clean-architecture/domain"
	"crud-clean-architecture/repository"
)

type CategoryService interface {
	CreateCategory(category *domain.Category) error
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id uint) error
	IsCategoryNameUnique(name string) (bool, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}

func (s *categoryService) CreateCategory(category *domain.Category) error {
	return s.categoryRepo.CreateCategory(category)
}

func (s *categoryService) IsCategoryNameUnique(name string) (bool, error) {
	return s.categoryRepo.IsCategoryNameUnique(name)
}
func (s *categoryService) GetAllCategories() ([]domain.Category, error) {
	return s.categoryRepo.GetAllCategories()
}

func (s *categoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	return s.categoryRepo.GetCategoryByID(id)
}

func (s *categoryService) UpdateCategory(category *domain.Category) error {
	return s.categoryRepo.UpdateCategory(category)
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.DeleteCategory(id)
}
