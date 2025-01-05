package repository

import (
	"context"
	"crud-clean-architecture/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *domain.Order) error
	GetAllOrders() ([]domain.Order, error)
	GetOrderByID(id uint) (*domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id uint) error
	CreateOrderWithDetails(order *domain.Order) error
	UpdateOrderInvoice(orderID uint, invoiceNumber string) error
}

type orderRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client) OrderRepository {
	return &orderRepository{db, redis}
}

const orderCacheKey = "order:all"

func (r *orderRepository) CreateOrderWithDetails(order *domain.Order) error {
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, orderCacheKey).Err(); err != nil {
		return err
	}
	// Mulai transaksi
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// Kosongkan field Details sebelum menyimpan order
	originalDetails := order.Details
	order.Details = nil

	// Simpan data order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Kembalikan data Details untuk diproses selanjutnya
	order.Details = originalDetails
	// Salin slice order.Details untuk menghindari konflik referensi
	details := make([]domain.OrderDetail, len(order.Details))
	copy(details, order.Details)
	// Simpan data order detail
	for i := range details {
		details[i].OrderID = order.ID                                       // Set OrderID untuk setiap detail
		fmt.Printf("Saving order detail #%d: %+v\n", i+1, order.Details[i]) // Logging

		if err := tx.Create(&details[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// Commit transaksi jika semua berhasil
	return tx.Commit().Error
}
func (r *orderRepository) CreateOrder(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetAllOrders() ([]domain.Order, error) {
	ctx := context.Background()

	// Cek cache
	cachedData, err := r.redis.Get(ctx, categoryCacheKey).Result()
	if err == nil {
		var orders []domain.Order
		if err := json.Unmarshal([]byte(cachedData), &orders); err == nil {
			return orders, nil
		}
	}

	// Jika cache tidak ada, fallback ke database
	var orders []domain.Order
	if err := r.db.Preload("Details").Preload("Details.Product").
		Preload("Details.Product.Category").Find(&orders).Error; err != nil {
		return nil, err
	}

	// Simpan ke cache
	data, _ := json.Marshal(orders)
	_ = r.redis.Set(ctx, orderCacheKey, data, 10*time.Minute).Err()
	return orders, nil

}

func (r *orderRepository) GetOrderByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Details").First(&order, id).Error
	return &order, err
}

func (r *orderRepository) UpdateOrder(order *domain.Order) error {
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, orderCacheKey).Err(); err != nil {
		return err
	}
	return r.db.Save(order).Error
}

func (r *orderRepository) DeleteOrder(id uint) error {
	ctx := context.Background()

	// Hapus cache setelah create
	if err := r.redis.Del(ctx, orderCacheKey).Err(); err != nil {
		return err
	}
	return r.db.Delete(&domain.Order{}, id).Error
}
func (r *orderRepository) UpdateOrderInvoice(orderID uint, invoiceNumber string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", orderID).Update("invoice_number", invoiceNumber).Error
}
