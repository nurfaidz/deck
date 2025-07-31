package controllers

import (
	"deck/services"
	"deck/structs"
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

	notification, total, err := nc.notificationService.GetNotifications(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch notifications",
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Notifications fetched successfully",
		Data: gin.H{
			"notifications": notification,
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
