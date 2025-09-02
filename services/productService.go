package services

import (
	"deck/enums"
	"deck/helpers"
	"deck/models"
	"deck/structs"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// Get All Products
func (ps *ProductService) GetProducts() ([]models.Product, error) {
	var products []models.Product
	err := ps.db.Order("created_at DESC").Find(&products).Error

	return products, err
}

// Get Product By Id
func (ps *ProductService) GetProductById(id uint) (*models.Product, error) {
	var product models.Product
	if err := ps.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (ps *ProductService) CreateProduct(req *structs.ProductCreateRequest, file *multipart.FileHeader) (*models.Product, error) {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if file == nil {
		tx.Rollback()
		return nil, errors.New("image file is required")
	}

	if file.Size > 1<<20 {
		tx.Rollback()
		return nil, errors.New("image size must be less than 1MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		tx.Rollback()
		return nil, errors.New("image must be a JPG, JPEG, or PNG file")
	}

	uploadDir := "uploads"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create upload directory: %v", err)
		}
	}

	if info, err := os.Stat(uploadDir); err != nil {
		log.Printf("ERROR: Failed to stat upload directory: %v", err)
	} else {
		log.Printf("INFO: Upload directory permissions: %s", info.Mode().Perm())
	}

	newFileName := helpers.GenerateUniqueFilename(file)
	imagePath := filepath.Join(uploadDir, newFileName)

	src, err := file.Open()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(imagePath)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save uploaded image: %v", err)
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
		tx.Rollback()
		return nil, errors.New("invalid category type")
	}

	product := models.Product{
		Name:        req.Name,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
		IsAvailable: req.IsAvailable,
		Image:       newFileName,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &product, nil
}

// Delete product
func (ps *ProductService) DeleteProduct(id uint) error {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var product models.Product
	if err := tx.First(&product, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to find product: %v", err)
	}

	if err := ps.db.Delete(&product).Error; err != nil {
		return err
	}

	// Delete image file
	imagePath := ""
	if product.Image != "" {
		imagePath = filepath.Join("uploads", product.Image)
	}

	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete product: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	if imagePath != "" {
		go func() {
			if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
				log.Printf("Warning: Failed to delete image %s: %v", imagePath, err)
			} else {
				log.Printf("Successfully deleted image: %s", imagePath)
			}
		}()
	}

	return nil

}

//func (ps *ProductService) UpdateProduct(id uint, req *structs.ProductUpdateRequest, file *multipart.FileHeader) (*models.Product, error) {
//	tx := ps.db.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			tx.Rollback()
//		}
//	}()
//
//	var product models.Product
//	if err := tx.First(&product, id).Error; err != nil {
//		tx.Rollback()
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, errors.New("product not found")
//		}
//		return nil, fmt.Errorf("failed to find product: %v", err)
//	}
//
//	if req.Category != "" {
//		validCategories := []enums.CategoryType{
//			enums.Classic, enums.Sparkling, enums.Smoothies,
//			enums.Tea, enums.IceCream, enums.Powders, enums.Other,
//		}
//
//		isValidCategory := false
//		for _, validCat := range validCategories {
//			if req.Category == validCat {
//				isValidCategory = true
//				break
//			}
//		}
//
//		if !isValidCategory {
//			tx.Rollback()
//			return nil, errors.New("invalid category type")
//		}
//	}
//
//	var newFileName string
//	var shouldUpdateImage bool
//	uploadDir := "uploads"
//
//	if file != nil {
//		if file.Size > 1<<20 { // 1MB
//			tx.Rollback()
//			return nil, errors.New("image size must be less than 1MB")
//		}
//
//		ext := strings.ToLower(filepath.Ext(file.Filename))
//		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
//			tx.Rollback()
//			return nil, errors.New("image must be a JPG, JPEG, or PNG file")
//		}
//
//		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
//			if err := os.MkdirAll(uploadDir, 0755); err != nil {
//				tx.Rollback()
//				return nil, fmt.Errorf("failed to create upload directory: %v", err)
//			}
//		}
//
//		newFileName = helpers.GenerateUniqueFilename(file)
//		imagePath := filepath.Join(uploadDir, newFileName)
//
//		src, err := file.Open()
//		if err != nil {
//			tx.Rollback()
//			return nil, fmt.Errorf("failed to open uploaded file: %v", err)
//		}
//		defer dst.Close()
//
//		dst, err := os.Create(imagePath)
//		if err != nil {
//			tx.Rollback()
//			return nil, fmt.Errorf("failed to create destination file: %v", err)
//		}
//		defer dst.Close()
//
//		if _, err := io.Copy(dst, src); err != nil {
//			tx.Rollback()
//			return nil, fmt.Errorf("failed to save uploaded image: %v", err)
//		}
//
//		shouldUpdateImage = true
//	}
//
//	oldImage := product.Image
//
//	product.Name = req.Name
//	product.Price = req.Price
//	product.Category = req.Category
//	product.Description = req.Description
//	product.IsAvailable = req.IsAvailable
//
//	if shouldUpdateImage {
//		product.Image = newFileName
//	}
//
//	if err := tx.Save(&product).Error; err != nil {
//		if shouldUpdateImage && newFileName != "" {
//			os.Remove(filepath.Join(uploadDir, newFileName))
//		}
//		tx.Rollback()
//		return nil, fmt.Errorf("failed to update product: %v", err)
//	}
//
//	if err := tx.Commit().Error; err != nil {
//		if shouldUpdateImage && newFileName != "" {
//			os.Remove(filepath.Join(uploadDir, newFileName))
//		}
//		return nil, fmt.Errorf("failed to commit transaction: %v", err)
//	}
//
//	if shouldUpdateImage && oldImage != "" {
//		go func() {
//			oldPath := filepath.Join(uploadDir, oldImage)
//			if err := os.Remove(oldPath); err != nil {
//				log.Printf("Warning: Failed to delete old image %s: %v", oldPath, err)
//			}
//		}()
//	}
//
//	return &product, nil
//}
