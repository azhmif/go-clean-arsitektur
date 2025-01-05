package repository

import (
	"crud-clean-architecture/domain"
	"errors"

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
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) CreateProduct(product *domain.Product) error {
	return r.db.Create(product).Error
}
func (r *productRepository) IsProductNameUnique(name string, categori_id uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Product{}).Where("name = ?", name).Where("category_id = ?", categori_id).Count(&count).Error
	return count == 0, err
}
func (r *productRepository) GetAllProducts() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
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

	// Hapus data jika ditemukan
	return r.db.Delete(&product).Error
}
