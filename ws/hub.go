package ws

import (
	"sync"
)

type Client struct {
	UserID uint
	Conn   *WebSocketConn
}

type WebSocketConn interface {
	WriteJSON(v interface{}) error
	Close() error
}

type Hub struct {
	clients map[uint][]WebSocketConn // userID -> conn listesi
	mu      sync.RWMutex
}

var NotificationHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint][]WebSocketConn),
	}
}

func (h *Hub) Register(userID uint, conn WebSocketConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], conn)
}

func (h *Hub) Unregister(userID uint, conn WebSocketConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[userID]
	for i, c := range conns {
		if c == conn {
			h.clients[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID)
	}
}

func (h *Hub) SendNotification(userID uint, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, conn := range h.clients[userID] {
		conn.WriteJSON(message)
	}
}
