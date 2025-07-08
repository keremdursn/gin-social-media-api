package models

import (
	"time"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}
