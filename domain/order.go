package domain

import "time"

type Order struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	InvoiceNumber string `json:"invoice_number"`

	OrderDate  time.Time     `json:"order_date"`
	TotalPrice float64       `json:"total_price"`
	Details    []OrderDetail `json:"details" gorm:"foreignKey:OrderID"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type OrderDetail struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
	Subtotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
