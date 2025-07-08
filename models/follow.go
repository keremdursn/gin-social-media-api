package models

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	FollowerID uint `json:"follower_id"` // Takip eden kullanıcı
	FollowingID uint `json:"following_id"` // Takip edilen kullanıcı
}
