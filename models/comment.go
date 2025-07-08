package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content" binding:"required"`
	UserID  uint   `json:"user_id"`
	PostID  uint   `json:"post_id"`
}