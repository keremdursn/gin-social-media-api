package controllers

import (
	"net/http"
	"strconv"

	"gin-blog-api/database"
	"gin-blog-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetCurrentUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş yapılmamış"})
		return
	}
	userID := uint(userIDVal.(float64))

	var user models.User
	if err := database.DB.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Şifreyi göstermiyoruz
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

func GetUserFollowCounts(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID"})
		return
	}

	var followersCount int64
	var followingCount int64

	db := database.DB.Db

	db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&followersCount)
	db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&followingCount)

	c.JSON(http.StatusOK, gin.H{
		"followers_count": followersCount,
		"following_count": followingCount,
	})
}

func GetUserProfile(c *gin.Context) {
	userIDParam := c.Param("id")
	profileUserID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID"})
		return
	}

	var user models.User
	if err := database.DB.Db.First(&user, profileUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Şifreyi gizle
	user.Password = ""

	// Giriş yapan kullanıcı ID'si (token'dan)
	currentUserIDVal, exists := c.Get("userID")
	var isFollowing bool
	if exists {
		currentUserID := uint(currentUserIDVal.(float64))

		var follow models.Follow
		err := database.DB.Db.Where("follower_id = ? AND following_id = ?", currentUserID, profileUserID).First(&follow).Error
		isFollowing = (err == nil)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"is_following": isFollowing,
	})
}

func GetFollowersCount(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID"})
		return
	}

	var count int64
	database.DB.Db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&count)

	c.JSON(http.StatusOK, gin.H{"followers_count": count})
}

func GetFollowingCount(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID"})
		return
	}

	var count int64
	database.DB.Db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&count)

	c.JSON(http.StatusOK, gin.H{"following_count": count})
}

func UpdateUserProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Bio      string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	var user models.User
	if err := database.DB.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	user.Username = input.Username
	user.Email = input.Email
	user.Bio = input.Bio

	if err := database.DB.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	user.Password = "" // şifreyi dışarı gösterme
	c.JSON(http.StatusOK, user)
}

func ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Eski ve yeni şifre gereklidir (yeni en az 6 karakter)"})
		return
	}

	var user models.User
	if err := database.DB.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Eski şifre kontrolü
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Eski şifre yanlış"})
		return
	}

	// Yeni şifre hash’le
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre hash’lenemedi"})
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Şifre başarıyla değiştirildi"})
}
