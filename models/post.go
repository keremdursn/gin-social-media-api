package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	UserID    uint      `json:"user_id"`
	Comments  []Comment `json:"comments" gorm:"foreignKey:PostID"`
	Likes     []Like    `json:"likes" gorm:"foreignKey:PostID"`
	ImageURLs []string  `gorm:"type:text[]" json:"image_urls"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
}
