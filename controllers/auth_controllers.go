package controllers

import (
	"gin-blog-api/database"
	"gin-blog-api/models"
	"gin-blog-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	db := database.DB.Db
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
		return
	}

	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre oluşturulamadı"})
		return
	}
	input.Password = string(hashedPassword)

	// Veritabanına kaydet

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı kaydedilemedi"})
		return
	}

	c.JSON(http.StatusCreated, input)

}

func Login(c *gin.Context) {

	db := database.DB.Db

	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
		return
	}

	// Kullanıcıyı email ile bul
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz email ya da şifre"})
		return
	}

	// Şifreyi karşılaştır
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz email ya da şifre"})
		return
	}

	// ✅ Token üretme işlemi utils üzerinden
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


