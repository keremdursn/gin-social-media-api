package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID      uint      `json:"user_id"`       // Bildirimi alan kullanıcı
	Content     string    `json:"content"`       // Bildirim metni
	IsRead      bool      `json:"is_read"`       // Okundu mu?
	CreatedAt   time.Time `json:"created_at"`
}
