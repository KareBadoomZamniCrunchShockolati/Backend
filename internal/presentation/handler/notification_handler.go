package handler

import (
	"net/http"
	"strconv"
	"time"

	dto "challenge-app/internal/application/dto"
	service "challenge-app/internal/application/service/interface"

	"github.com/gin-gonic/gin"
)

type NotificationHandlerImpl struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandlerImpl {
	return &NotificationHandlerImpl{svc: svc}
}

func (h *NotificationHandlerImpl) List(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userID := userIDAny.(uint)

	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	var cursor *time.Time
	if v := c.Query("cursor"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor"})
			return
		}
		cursor = &t
	}

	items, next, err := h.svc.ListRecent(c.Request.Context(), userID, limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load notifications"})
		return
	}

	resp := dto.NotificationListResponseDTO{
		Items:      make([]dto.NotificationResponseDTO, 0, len(items)),
		NextCursor: next,
	}

	for _, n := range items {
		resp.Items = append(resp.Items, dto.NotificationResponseDTO{
			ID:        n.ID,
			Type:      string(n.Type),
			TitleKey:  n.TitleKey,
			BodyKey:   n.BodyKey,
			Data:      n.Data,
			CreatedAt: n.CreatedAt,
			ReadAt:    n.ReadAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *NotificationHandlerImpl) MarkRead(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userID := userIDAny.(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.MarkRead(c.Request.Context(), userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark read"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *NotificationHandlerImpl) MarkAllRead(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userID := userIDAny.(uint)

	if err := h.svc.MarkAllRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark all read"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *NotificationHandlerImpl) UnreadCount(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userID := userIDAny.(uint)

	count, err := h.svc.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get unread count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
