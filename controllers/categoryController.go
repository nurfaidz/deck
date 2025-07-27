package controllers

import (
	"deck/enums"
	"deck/structs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCategories(c *gin.Context) {
	allCategories := enums.GetAllCategories()

	categories := make([]map[string]interface{}, len(allCategories))
	for i, category := range allCategories {
		categories[i] = map[string]interface{}{
			"value": string(category),
			"label": category.GetDisplayName(),
		}
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List of categories",
		Data:    categories,
	})
}
