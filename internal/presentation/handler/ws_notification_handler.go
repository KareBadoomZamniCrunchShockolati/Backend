package handler

import (
	"log"
	"net/http"
	"strings"

	"challenge-app/internal/infrastructure/realtime/ws"
	"challenge-app/pkg/security"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSNotificationHandlerImpl struct {
	hub      *ws.Hub
	jwt      security.JWTService
	upgrader websocket.Upgrader
}

func NewWSNotificationHandler(hub *ws.Hub, jwt security.JWTService) *WSNotificationHandlerImpl {
	return &WSNotificationHandlerImpl{
		hub: hub,
		jwt: jwt,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *WSNotificationHandlerImpl) Connect(c *gin.Context) {
	var token string

	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		token = strings.TrimPrefix(auth, "Bearer ")
	} else if q := c.Query("token"); q != "" {
		token = q
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := h.jwt.ValidateToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	userID := claims.UserID
	log.Println("WS connected user:", userID)
	wc := ws.NewConn(conn)
	h.hub.Add(userID, wc)

	wc.Send([]byte(`{"type":"connected","title":"WS Connected","body":"ok"}`))

	wc.ReadLoop(func() {
		h.hub.Remove(userID, wc)
	})
}
