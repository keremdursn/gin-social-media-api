package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gin-blog-api/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token gerekli"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Sunucu yapılandırması hatası"})
			c.Abort()
			return
		}

		// JWT token doğrulaması
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// HMAC metodu kullanıldığını doğrula
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			c.Abort()
			return
		}

		// Redis session kontrolü
		_, err = utils.GetSession(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session bulunamadı veya süresi dolmuş"})
			c.Abort()
			return
		}

		// Token içindeki claim'lerden userID al
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token çözümlenemedi"})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token içinde kullanıcı bilgisi yok"})
			c.Abort()
			return
		}

		// Context'e userID'yi set et (string olarak, istersen uint dönüştür)
		c.Set("userID", userID)
		c.Next()
	}
}
