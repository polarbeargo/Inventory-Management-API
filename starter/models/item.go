package models

type Item struct {
	ID    string  `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name" gorm:"not null" binding:"required,min=1,max=100"`
	Stock int     `json:"stock" gorm:"not null" binding:"required,min=0"`
	Price float64 `json:"price" gorm:"not null" binding:"required,gt=0"`
}
