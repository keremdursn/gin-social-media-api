package controllers

import (
	"log"
	"net/http"
	"strconv"

	"gin-blog-api/database"
	"gin-blog-api/models"
	"gin-blog-api/utils"

	"github.com/gin-gonic/gin"
)

func ToggleLike(c *gin.Context) {
	// Giriş yapmış kullanıcı ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	// Post ID'yi parametreden al
	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz gönderi ID"})
		return
	}

	var existingLike models.Like
	result := database.DB.Db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike)

	if result.Error == nil {
		// Zaten beğenmişse, kaldır (unlike)
		if err := database.DB.Db.Delete(&existingLike).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Beğeni kaldırılamadı"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Beğeni kaldırıldı"})
		return
	}

	// Beğeni yoksa, oluştur
	newLike := models.Like{
		UserID: userID,
		PostID: uint(postID),
	}

	if err := database.DB.Db.Create(&newLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Beğeni eklenemedi"})
		return
	}

	if err := utils.CreateLikeNotification(userID, newLike.UserID, newLike.ID); err != nil {
		log.Println("Bildirim oluşturulamadı:", err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Gönderi beğenildi"})
}
