package domain

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null"`
}

type CategoryForm struct {
	Name string `json:"name" binding:"required,max=255"`
}
