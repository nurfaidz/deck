package controllers

import (
	"deck/database"
	"deck/enums"
	"deck/helpers"
	"deck/models"
	"deck/structs"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetProducts(c *gin.Context) {
	name := strings.TrimSpace(c.Query("filter[name]"))
	category := strings.TrimSpace(c.Query("filter[category]"))

	var products []models.Product
	query := database.DB

	if name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch products",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	var filters []string
	if name != "" {
		filters = append(filters, fmt.Sprintf("Name: %s", name))
	}

	if category != "" {
		filters = append(filters, fmt.Sprintf("category: %s", category))
	}

	message := "All products"
	if len(filters) > 0 {
		message = fmt.Sprintf("Products filteres by %s", strings.Join(filters, ", "))
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data:    products,
	})
}

func CreateProduct(c *gin.Context) {
	var req = structs.ProductCreateRequest{}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if file.Size > 1<<20 { // 1MB
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"image": "Image size must be less than 1MB"},
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"image": "Image must be a JPG, JPEG, or PNG file"},
		})

		return
	}

	uploadDir := "uploads"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Failed to create upload directory",
				Errors:  map[string]string{"image": "Failed to create upload directory"},
			})

			return
		}
	}

	// check directory permissions
	if info, err := os.Stat(uploadDir); err != nil {
		log.Printf("ERROR: Failed to stat upload directory: %v", err)
	} else {
		log.Printf("INFO: Upload directory permissions: %s", info.Mode().Perm())
	}

	newFileName := helpers.GenerateUniqueFilename(file)
	imagePath := filepath.Join(uploadDir, newFileName)

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to upload image",
			Errors:  map[string]string{"image": "Failed to save uploaded image"},
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
		})

		return
	}

	product := models.Product{
		Name:        req.Name,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
		IsAvailable: req.IsAvailable,
		Image:       newFileName,
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

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Validate category if provided
	if req.Category != "" {
		validCategories := []enums.CategoryType{
			enums.Classic, enums.Sparkling, enums.Smoothies,
			enums.Tea, enums.IceCream, enums.Powders, enums.Other,
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
			})
			return
		}
	}

	var newFileName string
	var shouldUpdateImage bool

	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > 1<<20 { // 1MB
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"image": "Image size must be less than 1MB"},
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"image": "Image must be a JPG, JPEG, or PNG file"},
			})
			return
		}

		newFileName = helpers.GenerateUniqueFilename(file)
		uploadDir := "uploads"
		imagePath := filepath.Join(uploadDir, newFileName)

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Failed to upload image",
				Errors:  map[string]string{"image": "Failed to save uploaded image"},
			})
			return
		}

		shouldUpdateImage = true
	} else if err.Error() != "http: no such file" {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Image upload error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Store old image for cleanup
	oldImage := product.Image

	product.Name = req.Name
	product.Price = req.Price
	product.Category = req.Category
	product.Description = req.Description
	product.IsAvailable = req.IsAvailable

	if shouldUpdateImage {
		product.Image = newFileName
	}

	if err := database.DB.Save(&product).Error; err != nil {
		if shouldUpdateImage {
			os.Remove(filepath.Join("uploads", newFileName))
		}

		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update product",
		})
		return
	}

	if shouldUpdateImage && oldImage != "" {
		go func() {
			oldPath := filepath.Join("uploads", oldImage)
			if err := os.Remove(oldPath); err != nil {
				log.Printf("Warning: Failed to delete old image %s: %v", oldPath, err)
			}
		}()
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

func FilterByName(c *gin.Context) {
	name := c.Query("name")
	var products []models.Product

	if name != "" {
		if err := database.DB.Where("LOWER(name) LIKE LOWER(?)", "%"+strings.TrimSpace(name)+"%").Find(&products).Error; err != nil {
			c.JSON(http.StatusNotFound, structs.ErrorResponse{
				Success: false,
				Message: "No products found",
				Errors:  helpers.TranslateErrorMessage(err),
			})

			return
		}
	} else {
		database.DB.Find(&products)
	}

	message := "All products"
	if name != "" {
		message = fmt.Sprintf("Filtered products by name: %s", name)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data:    products,
	})
}

func FilterByCategory(c *gin.Context) {
	category := c.Query("category")
	var products []models.Product

	query := database.DB

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch products",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	message := "All products"
	if category != "" {
		message = fmt.Sprintf("Products filtered by category: %s", category)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data:    products,
	})
}
