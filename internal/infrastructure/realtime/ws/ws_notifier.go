package ws

import (
	"context"
	"encoding/json"

	appsvc "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/model"
)

type WSNotifierImpl struct {
	hub *Hub
}

func NewWSNotifier(hub *Hub) appsvc.Notifier {
	return &WSNotifierImpl{hub: hub}
}

type notificationWire struct {
	ID        uint                   `json:"id"`
	Type      model.NotificationType `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Data      map[string]any         `json:"data,omitempty"`
	CreatedAt string                 `json:"created_at"` // client friendly (RFC3339)
	ReadAt    *string                `json:"read_at,omitempty"`
}

func (n *WSNotifierImpl) Push(ctx context.Context, userID uint, notif model.Notification) error {
	created := notif.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	var readAt *string
	if notif.ReadAt != nil {
		s := notif.ReadAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		readAt = &s
	}

	wire := notificationWire{
		ID:        notif.ID,
		Type:      notif.Type,
		Title:     notif.Title,
		Body:      notif.Body,
		Data:      notif.Data,
		CreatedAt: created,
		ReadAt:    readAt,
	}

	b, _ := json.Marshal(wire)
	n.hub.Send(userID, b)
	return nil
}
