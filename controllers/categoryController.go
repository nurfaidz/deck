package controllers

import (
	"deck/services"
	"deck/structs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (cc *CategoryController) GetCategories(c *gin.Context) {
	categories := cc.categoryService.GetCategories()

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List of categories",
		Data:    categories,
	})
}

func (cc *CategoryController) GetCategoryByValue(c *gin.Context) {
	value := c.Param("value")

	if value == "" {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Category value is required",
		})

		return
	}

	category, err := cc.categoryService.GetCategoryByValue(value)
	if err != nil {
		var statusCode int
		var message string

		if strings.Contains(err.Error(), "invalid category value") {
			statusCode = http.StatusNotFound
			message = "Category not found"
		} else {
			statusCode = http.StatusInternalServerError
			message = "Internal server error"
		}

		c.JSON(statusCode, structs.ErrorResponse{
			Success: false,
			Message: message,
			Errors:  map[string]string{"error": err.Error()},
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category found",
		Data:    category,
	})
}
