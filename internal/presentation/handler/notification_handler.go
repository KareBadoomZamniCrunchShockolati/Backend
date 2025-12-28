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

// GET /api/v1/notifications?limit=20&cursor=2025-12-27T10:00:00Z
func (h *NotificationHandlerImpl) List(c *gin.Context) {
	userID := c.GetUint("user_id")

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
			Title:     n.Title,
			Body:      n.Body,
			Data:      n.Data,
			CreatedAt: n.CreatedAt,
			ReadAt:    n.ReadAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// POST /api/v1/notifications/:id/read
func (h *NotificationHandlerImpl) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")

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

// POST /api/v1/notifications/read-all
func (h *NotificationHandlerImpl) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.svc.MarkAllRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark all read"})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /api/v1/notifications/unread-count
func (h *NotificationHandlerImpl) UnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.svc.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get unread count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
