package service

import (
	"crud-clean-architecture/domain"
	"crud-clean-architecture/repository"
	"errors"
	"fmt"
	"time"
)

type OrderService interface {
	CreateOrder(order *domain.Order) error
	GetAllOrders() ([]domain.Order, error)
	GetOrderByID(id uint) (*domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id uint) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{orderRepo, productRepo}
}

func (s *orderService) CreateOrder(order *domain.Order) error {
	var totalPrice float64
	for i := range order.Details {
		product, err := s.productRepo.GetProductByID(order.Details[i].ProductID)
		if err != nil {
			return errors.New("product not found")
		}
		order.Details[i].Subtotal = product.Price * float64(order.Details[i].Quantity)
		totalPrice += order.Details[i].Subtotal
	}

	order.TotalPrice = totalPrice
	order.OrderDate = time.Now()

	// Simpan order untuk mendapatkan ID
	err := s.orderRepo.CreateOrderWithDetails(order)
	if err != nil {
		return err
	}

	// Generate kode invoice
	invoiceNumber := s.generateInvoiceNumber(order.ID, order.OrderDate)

	// Perbarui invoice number di database
	return s.orderRepo.UpdateOrderInvoice(order.ID, invoiceNumber)
}

func (s *orderService) GetAllOrders() ([]domain.Order, error) {
	return s.orderRepo.GetAllOrders()
}

func (s *orderService) GetOrderByID(id uint) (*domain.Order, error) {
	return s.orderRepo.GetOrderByID(id)
}

func (s *orderService) UpdateOrder(order *domain.Order) error {
	return s.orderRepo.UpdateOrder(order)
}

func (s *orderService) DeleteOrder(id uint) error {
	return s.orderRepo.DeleteOrder(id)
}
func (s *orderService) generateInvoiceNumber(orderID uint, orderDate time.Time) string {
	datePart := orderDate.Format("20060102") // Format tanggal menjadi YYYYMMDD
	return fmt.Sprintf("INV-%s-%d", datePart, orderID)
}
