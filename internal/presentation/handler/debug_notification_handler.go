package handler

import (
	"net/http"

	service "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type DebugHandlerImpl struct {
	notifSvc service.NotificationService
}

func NewDebugHandler(notifSvc service.NotificationService) *DebugHandlerImpl {
	return &DebugHandlerImpl{
		notifSvc: notifSvc,
	}
}

func (h *DebugHandlerImpl) SendTestNotification(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	err := h.notifSvc.CreateAndPush(
		c.Request.Context(),
		model.Notification{
			UserID:   userID,
			Type:     model.NotifFollowed, 
			TitleKey: "notif.debug.title",
			BodyKey:  "notif.debug.body",
			Data: map[string]any{
				"ok":       true,
				"user_id":  userID,
				"scenario": "debug",
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sent",
		"user_id": userID,
	})
}
