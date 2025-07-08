package utils

import (
	"fmt"
	"gin-blog-api/database"
	"gin-blog-api/models"
)

func CreateFollowNotification(followerID, followingID uint) error {
	content := fmt.Sprintf("Yeni bir takipçiniz var! Kullanıcı ID: %d", followerID)
	notification := models.Notification{
		UserID:  followingID,
		Content: content,
		IsRead:  false,
	}
	return database.DB.Db.Create(&notification).Error
}


// Gönderiye yorum yapılınca
func CreateCommentNotification(commenterID, postOwnerID, postID uint) error {
	if commenterID == postOwnerID {
		return nil // kendi postuna yorum yaptıysa bildirim oluşturma
	}
	content := fmt.Sprintf("Gönderinize bir yorum yapıldı! (Post ID: %d)", postID)
	notification := models.Notification{
		UserID:  postOwnerID,
		Content: content,
		IsRead:  false,
	}
	return database.DB.Db.Create(&notification).Error
}

// Gönderi beğenilince
func CreateLikeNotification(likerID, postOwnerID, postID uint) error {
	if likerID == postOwnerID {
		return nil // kendi postunu beğendiyse bildirim oluşturma
	}
	content := fmt.Sprintf("Gönderinize bir beğeni geldi! (Post ID: %d)", postID)
	notification := models.Notification{
		UserID:  postOwnerID,
		Content: content,
		IsRead:  false,
	}
	return database.DB.Db.Create(&notification).Error
}