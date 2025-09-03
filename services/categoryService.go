package services

import (
	"deck/enums"
	"deck/structs"
	"errors"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

// Get All Categories
func (cs *CategoryService) GetCategories() []structs.CategoryResponse {
	allCategories := enums.GetAllCategories()

	categories := make([]structs.CategoryResponse, len(allCategories))
	for i, category := range allCategories {
		categories[i] = structs.CategoryResponse{
			Value: string(category),
			Label: category.GetDisplayName(),
		}
	}

	return categories
}

// Get Category By Value returns a single category by its value
func (cs *CategoryService) GetCategoryByValue(value string) (*structs.CategoryResponse, error) {
	categoryType := enums.CategoryType(value)

	allCategories := enums.GetAllCategories()
	isValid := false

	for _, validCategory := range allCategories {
		if validCategory == categoryType {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("invalid category value")
	}

	return &structs.CategoryResponse{
		Value: string(categoryType),
		Label: categoryType.GetDisplayName(),
	}, nil
}
