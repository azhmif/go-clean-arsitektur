package repository

import (
	"crud-clean-architecture/domain"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrCategoryNameExists = errors.New("category name already exists")
)

type CategoryRepository interface {
	CreateCategory(category *domain.Category) error
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id uint) error
	IsCategoryNameUnique(name string) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) CreateCategory(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) IsCategoryNameUnique(name string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Category{}).Where("name = ?", name).Count(&count).Error
	return count == 0, err
}
func (r *categoryRepository) GetAllCategories() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetCategoryByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) UpdateCategory(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	var category domain.Category

	// Periksa apakah data dengan ID ada
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Hapus data jika ditemukan
	return r.db.Delete(&category).Error
}
