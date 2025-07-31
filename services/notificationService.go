package services

import (
	"deck/models"
	"encoding/json"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// Broadcasts notification to all admin
func (ns *NotificationService) BroadcastToAdmins(notificationType, title, message, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserId:  nil,
		Type:    notificationType.(string),
		Title:   title.(string),
		Message: message.(string),
		Data:    string(dataBytes),
		IsRead:  false,
	}
	return ns.db.Create(notification).Error
}

// Get notification for all admin
func (ns *NotificationService) GetNotifications(Id uint) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	// query notifications for spesific admin
	query := ns.db.Where("user_id = ? OR user_id IS NULL", Id)

	// count notifications
	query.Model(&models.Notification{}).Count(&total)

	err := query.Order("created_at DESC").Find(&notifications).Error

	return notifications, total, err
}

// Get unread notifications count
func (ns *NotificationService) GetUnreadNotificationCount(Id uint) (int64, error) {
	var count int64
	err := ns.db.Model(&models.Notification{}).Where("(user_id = ? OR user_id IS NULL) AND is_read = false", Id).Count(&count).Error

	return count, err
}

// Mark notification as read
func (ns *NotificationService) MarkAsRead(notificationId uint, Id uint) error {
	return ns.db.Model(&models.Notification{}).Where("id = ? AND (user_id = ? OR user_id IS NULL)", notificationId, Id).Update("is_read", true).Error
}

// Mark all notifications as read
func (ns *NotificationService) MarkAllAsRead(Id uint) error {
	return ns.db.Model(&models.Notification{}).
		Where("(user_id = ? OR user_id IS NULL) AND is_read = false", Id).
		Update("is_read", true).Error
}
