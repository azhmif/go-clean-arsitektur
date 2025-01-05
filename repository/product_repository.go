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

type ProductRepository interface {
	CreateProduct(product *domain.Product) error
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id uint) error
	IsProductNameUnique(name string, categori_id uint) (bool, error)
}

type productRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepository(db *gorm.DB, redis *redis.Client) ProductRepository {
	return &productRepository{db, redis}
}

const productCacheKey = "product:all"

func (r *productRepository) CreateProduct(product *domain.Product) error {
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, productCacheKey).Err(); err != nil {
		return err
	}
	return r.db.Create(product).Error
}
func (r *productRepository) IsProductNameUnique(name string, categori_id uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Product{}).Where("name = ?", name).Where("category_id = ?", categori_id).Count(&count).Error
	return count == 0, err
}
func (r *productRepository) GetAllProducts() ([]domain.Product, error) {
	ctx := context.Background()

	// Cek cache
	cachedData, err := r.redis.Get(ctx, productCacheKey).Result()
	if err == nil {
		var products []domain.Product
		if err := json.Unmarshal([]byte(cachedData), &products); err == nil {
			return products, nil
		}
	}

	// Jika cache tidak ada, fallback ke database
	var products []domain.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}

	// Simpan ke cache
	data, _ := json.Marshal(products)
	_ = r.redis.Set(ctx, productCacheKey, data, 10*time.Minute).Err()

	return products, nil
}

func (r *productRepository) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateProduct(product *domain.Product) error {
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, productCacheKey).Err(); err != nil {
		return err
	}
	return r.db.Save(product).Error
}

func (r *productRepository) DeleteProduct(id uint) error {
	var product domain.Product
	// Periksa apakah data dengan ID ada
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, productCacheKey).Err(); err != nil {
		return err
	}
	// Hapus data jika ditemukan
	return r.db.Delete(&product).Error
}
