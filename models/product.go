package models

import "deck/enums"

type Product struct {
	GormModel
	Name        string             `json:"name" gorm:"not null"`
	Category    enums.CategoryType `json:"category" gorm:"not null"`
	Description string             `json:"description" gorm:"type:text"`
	Image       string             `json:"image" gorm:"type:varchar(255)"`
	Price       uint               `json:"price" gorm:"not null"`
	IsAvailable bool               `json:"is_available" gorm:"not null;default:true"`
}
