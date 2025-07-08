package controllers

import (
	"deck/database"
	"deck/enums"
	"deck/helpers"
	"deck/models"
	"deck/structs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	database.DB.Find(&products)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List products data",
		Data:    products,
	})
}

func CreateProduct(c *gin.Context) {
	var req = structs.ProductCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	validCategories := []enums.CategoryType{
		enums.Classic, enums.Sparkling, enums.Smoothies, enums.Tea, enums.IceCream, enums.Powders, enums.Other,
	}

	isValidCategory := false
	for _, validCat := range validCategories {
		if req.Category == validCat {
			isValidCategory = true
			break
		}
	}

	if !isValidCategory {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Invalid category type",
			Errors:  map[string]string{"category": "Category must be one of the following: classic, sparkling, smoothies, tea, powders, ice_cream, other"},
		})

		return
	}

	product := models.Product{
		Name:        req.Name,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
		IsAvailable: req.IsAvailable,
		Image:       req.Image,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create product",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Product created successfully",
		Data: structs.ProductResponse{
			Id:           product.Id,
			Name:         product.Name,
			Price:        product.Price,
			Category:     product.Category,
			CategoryName: product.Category.GetDisplayName(),
			Image:        product.Image,
			Description:  product.Description,
			IsAvailable:  product.IsAvailable,
			CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func GetProductById(c *gin.Context) {
	var product models.Product
	productId := c.Param("id")

	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Product found",
		Data: structs.ProductResponse{
			Id:           product.Id,
			Name:         product.Name,
			Price:        product.Price,
			Category:     product.Category,
			CategoryName: product.Category.GetDisplayName(),
			Image:        product.Image,
			Description:  product.Description,
			IsAvailable:  product.IsAvailable,
			CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	productId := c.Param("id")

	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	var req = structs.ProductUpdateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	if req.Category != "" {
		validCategories := []enums.CategoryType{
			enums.Classic, enums.Sparkling, enums.Smoothies, enums.Tea, enums.IceCream, enums.Powders, enums.Other,
		}

		isValidCategory := false
		for _, validCat := range validCategories {
			if req.Category == validCat {
				isValidCategory = true
				break
			}
		}

		if !isValidCategory {
			c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
				Success: false,
				Message: "Invalid category type",
				Errors:  map[string]string{"category": "Category must be one of the following: classic, sparkling, smoothies, tea, powders, ice_cream, other"},
			})
			return
		}
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Category = req.Category
	product.Description = req.Description
	product.Image = req.Image
	product.IsAvailable = req.IsAvailable

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update product",
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Product updated successfully",
		Data: structs.ProductResponse{
			Id:           product.Id,
			Name:         product.Name,
			Price:        product.Price,
			Category:     product.Category,
			CategoryName: product.Category.GetDisplayName(),
			Image:        product.Image,
			Description:  product.Description,
			IsAvailable:  product.IsAvailable,
			CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func DeleteProduct(c *gin.Context) {
	var product models.Product
	productId := c.Param("id")

	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete product",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}
