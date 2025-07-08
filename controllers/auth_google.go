package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gin-blog-api/database"
	"gin-blog-api/models"
	"gin-blog-api/oauth"
	"gin-blog-api/utils"
)

func GoogleLogin(c *gin.Context) {
	// ileride CSRF için state kontrolü ekleyebilirsin
	url := oauth.GetGoogleLoginURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := oauth.ExchangeCodeForToken(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token alınamadı"})
		return
	}

	client := oauth.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı bilgisi alınamadı"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var userInfo struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Id       string `json:"id"`
		Verified bool   `json:"verified_email"`
	}
	json.Unmarshal(body, &userInfo)

	// Kullanıcıyı veritabanında kontrol et
	var user models.User
	result := database.DB.Db.Where("email = ?", userInfo.Email).First(&user)

	if result.Error != nil || user.ID == 0 {
		// Kullanıcı yoksa kayıt et
		user = models.User{
			Email:    userInfo.Email,
			Username: userInfo.Name,
			Bio:      "Google ile kayıt oldu.",
		}
		database.DB.Db.Create(&user)
	}

	// JWT + Redis session oluştur
	tokenString, err := utils.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Google login başarılı",
		"token":   tokenString,
	})
}
