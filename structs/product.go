package structs

import "deck/enums"

type ProductResponse struct {
	Id           uint               `json:"id"`
	Name         string             `json:"name"`
	Price        uint               `json:"price"`
	Category     enums.CategoryType `json:"category"`
	CategoryName string             `json:"category_name"`
	Description  string             `json:"description"`
	Image        string             `json:"image"`
	IsAvailable  bool               `json:"is_available"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}

type ProductCreateRequest struct {
	Name        string             `form:"name" binding:"required" gorm:"not null"`
	Price       uint               `form:"price" binding:"required" gorm:"not null"`
	Category    enums.CategoryType `form:"category" binding:"required,oneof=classic sparkling smoothies tea powders ice_cream other" gorm:"not null"`
	Description string             `form:"description"`
	IsAvailable bool               `form:"is_available" binding:"required" gorm:"not null"`
}

type ProductUpdateRequest struct {
	Name        string             `form:"name" binding:"required" gorm:"not null"`
	Price       uint               `form:"price" binding:"required" gorm:"not null"`
	Category    enums.CategoryType `form:"category" binding:"required" gorm:"not null"`
	Description string             `form:"description"`
	IsAvailable bool               `form:"is_available"`
}
