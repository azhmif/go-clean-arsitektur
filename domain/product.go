package domain

type Product struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type ProductForm struct {
	Name       string  `json:"name" binding:"required,max=255"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	CategoryID uint    `json:"category_id" binding:"required"`
}
