package controllers

import (
	"net/http"

	"gin-blog-api/database"
	"gin-blog-api/models"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	var notifications []models.Notification
	if err := database.DB.Db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bildirimler alınamadı"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Db.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Okundu olarak işaretlenemedi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bildirim okundu"})
}
