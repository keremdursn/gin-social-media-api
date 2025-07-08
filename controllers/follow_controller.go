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

func FollowUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	followerID := uint(userIDVal.(float64))

	followingIDParam := c.Param("id")
	followingID, err := strconv.Atoi(followingIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID"})
		return
	}

	if followerID == uint(followingID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kendini takip edemezsin"})
		return
	}

	// Takip zaten varsa sil (unfollow)
	var existing models.Follow
	result := database.DB.Db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existing)

	if result.Error == nil {
		if err := database.DB.Db.Delete(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Takipten çıkılamadı"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Takipten çıktınız"})
		return
	}

	// Yeni takip oluştur
	follow := models.Follow{
		FollowerID:  followerID,
		FollowingID: uint(followingID),
	}

	if err := database.DB.Db.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Takip yapılamadı"})
		return
	}

	// Takip oluşturduktan sonra
	if err := utils.CreateFollowNotification(followerID, uint(followingID)); err != nil {
		log.Println("Bildirim oluşturulamadı:", err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kullanıcı takip edildi"})
}
