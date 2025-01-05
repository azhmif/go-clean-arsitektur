package service

import (
	"crud-clean-architecture/domain"
	"crud-clean-architecture/repository"
)

type ProductService interface {
	CreateProduct(product *domain.Product) error
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id uint) error
	IsProductNameUnique(name string, categori_id uint) (bool, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo}
}

func (s *productService) CreateProduct(product *domain.Product) error {
	return s.productRepo.CreateProduct(product)
}

func (s *productService) GetAllProducts() ([]domain.Product, error) {
	return s.productRepo.GetAllProducts()
}

func (s *productService) GetProductByID(id uint) (*domain.Product, error) {
	return s.productRepo.GetProductByID(id)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	return s.productRepo.UpdateProduct(product)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.productRepo.DeleteProduct(id)
}

func (s *productService) IsProductNameUnique(name string, categori_id uint) (bool, error) {
	return s.productRepo.IsProductNameUnique(name, categori_id)
}
