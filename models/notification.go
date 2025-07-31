package models

type Notification struct {
	GormModel
	UserId  *uint  `json:"user_id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Data    string `json:"data"`
	IsRead  bool   `json:"is_read" gorm:"default:false"`
}
