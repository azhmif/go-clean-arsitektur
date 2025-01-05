package repository

import (
	"context"
	"crud-clean-architecture/domain"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
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
	db    *gorm.DB
	redis *redis.Client
}

func NewCategoryRepository(db *gorm.DB, redis *redis.Client) CategoryRepository {
	return &categoryRepository{db, redis}
}

const categoryCacheKey = "categories:all"

func (r *categoryRepository) CreateCategory(category *domain.Category) error {
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, categoryCacheKey).Err(); err != nil {
		return err
	}
	return r.db.Create(category).Error
}

func (r *categoryRepository) IsCategoryNameUnique(name string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Category{}).Where("name = ?", name).Count(&count).Error
	return count == 0, err
}
func (r *categoryRepository) GetAllCategories() ([]domain.Category, error) {
	ctx := context.Background()

	// Cek cache
	cachedData, err := r.redis.Get(ctx, categoryCacheKey).Result()
	if err == nil {
		var categories []domain.Category
		if err := json.Unmarshal([]byte(cachedData), &categories); err == nil {
			return categories, nil
		}
	}

	// Jika cache tidak ada, fallback ke database
	var categories []domain.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	// Simpan ke cache
	data, _ := json.Marshal(categories)
	_ = r.redis.Set(ctx, categoryCacheKey, data, 10*time.Minute).Err()

	return categories, nil
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
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, categoryCacheKey).Err(); err != nil {
		return err
	}
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
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, categoryCacheKey).Err(); err != nil {
		return err
	}
	// Hapus data jika ditemukan
	return r.db.Delete(&category).Error
}
