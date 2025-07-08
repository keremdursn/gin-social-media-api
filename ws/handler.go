package ws

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // güvenlik için ileride kontrol eklersin
	},
}

func NotificationSocket(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id gerekli"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz user_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	NotificationHub.Register(uint(userID), conn)
	defer NotificationHub.Unregister(uint(userID), conn)

	for {
		// WebSocket bağlantısı açık kalsın diye boş dinleme
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
