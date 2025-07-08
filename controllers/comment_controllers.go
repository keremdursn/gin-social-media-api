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

func CreateComment(c *gin.Context) {
	db := database.DB.Db

	// Kullanıcı ID'sini al
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	// Post ID'yi URL'den al
	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz gönderi ID"})
		return
	}

	// Yorum içeriğini al
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yorum içeriği zorunlu"})
		return
	}

	// Set UserID ve PostID
	comment.UserID = userID
	comment.PostID = uint(postID)

	// Kaydet
	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorum kaydedilemedi"})
		return
	}

	if err := utils.CreateCommentNotification(userID, comment.UserID, comment.ID); err != nil {
		log.Println("Bildirim oluşturulamadı:", err)
	}

	c.JSON(http.StatusCreated, comment)
}

func UpdateComment(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	commentID := c.Param("id")
	var comment models.Comment

	if err := database.DB.Db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Yorum bulunamadı"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu yorumu güncelleyemezsiniz"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "İçerik zorunlu"})
		return
	}

	comment.Content = input.Content

	if err := database.DB.Db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorum güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	commentID := c.Param("id")
	var comment models.Comment

	if err := database.DB.Db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Yorum bulunamadı"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu yorumu silemezsiniz"})
		return
	}

	if err := database.DB.Db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorum silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Yorum silindi"})
}
