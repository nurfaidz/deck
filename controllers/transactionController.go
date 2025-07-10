package controllers

import (
	"deck/database"
	"deck/models"
	"deck/structs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	database.DB.Find(&transactions).Preload("TransactionDetails")

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List transactions data",
		Data:    transactions,
	})
}
