package realtime

import (
	"net/http"

	"manumental-effort/server/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	hub          *Hub
	tokenManager *auth.TokenManager
	upgrader     websocket.Upgrader
}

func NewHandler(hub *Hub, tokenManager *auth.TokenManager) *Handler {
	return &Handler{
		hub:          hub,
		tokenManager: tokenManager,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) ServeWS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	userIDHex, err := h.tokenManager.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user id"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(h.hub, conn, userID)
	h.hub.Register(client)

	go client.WritePump()
	client.ReadPump()
}
