package controllers

import (
	"deck/services"
	"deck/structs"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

// Get notifications
func (nc *NotificationController) GetNotifications(c *gin.Context) {
	Id := c.GetUint("user_id")

	notifications, total, err := nc.notificationService.GetNotifications(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch notifications",
		})

		return
	}

	var result []map[string]interface{}
	for _, notification := range notifications {
		notifMap := map[string]interface{}{
			"id":         notification.Id,
			"created_at": notification.CreatedAt,
			"updated_at": notification.UpdatedAt,
			"user_id":    notification.UserId,
			"type":       notification.Type,
			"title":      notification.Title,
			"message":    notification.Message,
			"is_read":    notification.IsRead,
		}

		var dataObj map[string]interface{}
		if notification.Data != "" {
			json.Unmarshal([]byte(notification.Data), &dataObj)

			if products, exists := dataObj["products"]; exists {
				if productsArray, ok := products.([]interface{}); ok {
					var cleanProducts []map[string]interface{}
					for _, product := range productsArray {
						if productMap, ok := product.(map[string]interface{}); ok {
							cleanProduct := make(map[string]interface{})
							for key, value := range productMap {
								if key != "product" && key != "transaction" && value != nil {
									if str, ok := value.(string); ok && str != "" {
										cleanProduct[key] = value
									} else if num, ok := value.(float64); ok {
										cleanProduct[key] = num
									} else if _, ok := value.(bool); ok {
										cleanProduct[key] = value
									}
								}
							}

							if len(cleanProduct) > 0 {
								cleanProducts = append(cleanProducts, cleanProduct)
							}
						}
					}

					transaction := map[string]interface{}{
						"id":                  dataObj["transaction_id"],
						"transaction_details": cleanProducts,
					}

					dataObj["transaction"] = transaction
					delete(dataObj, "products")
				}
			}

			cleanDataObj := make(map[string]interface{})
			for key, value := range dataObj {
				if value != nil {
					if str, ok := value.(string); ok && str != "" {
						cleanDataObj[key] = value
					} else if _, ok := value.(float64); ok {
						cleanDataObj[key] = value
					} else if _, ok := value.(map[string]interface{}); ok {
						cleanDataObj[key] = value
					}
				}
			}
			dataObj = cleanDataObj
		}

		notifMap["data"] = dataObj

		result = append(result, notifMap)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Notifications fetched successfully",
		Data: gin.H{
			"notifications": result,
			"total":         total,
		},
	})
}

// Get unread notifications count
func (nc *NotificationController) GetUnreadCount(c *gin.Context) {
	Id := c.GetUint("user_id")

	count, err := nc.notificationService.GetUnreadNotificationCount(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch unread notifications count",
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Data: gin.H{
			"unread_count": count,
		},
	})
}

// Mark notification as read
func (nc *NotificationController) MarkAsRead(c *gin.Context) {
	Id := c.GetUint("user_id")
	notificationId, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := nc.notificationService.MarkAsRead(uint(notificationId), Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to mark notification as read",
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Notification marked as read",
	})
}

// Mark all notifications as read
func (nc *NotificationController) MarkAllAsRead(c *gin.Context) {
	Id := c.GetUint("user_id")

	err := nc.notificationService.MarkAllAsRead(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to mark all notifications as read",
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "All notifications marked as read",
	})
}
