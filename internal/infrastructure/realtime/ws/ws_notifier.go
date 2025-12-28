package ws

import (
	"context"
	"encoding/json"
	"time"

	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/model"
)

type WSNotifier struct {
	hub *Hub
}

func NewWSNotifier(hub *Hub) serviceinterface.Notifier {
	return &WSNotifier{hub: hub}
}

type notificationWire struct {
	ID        uint                   `json:"id"`
	Type      model.NotificationType `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Data      map[string]any         `json:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	ReadAt    *time.Time             `json:"read_at,omitempty"`
}

func (n *WSNotifier) Push(ctx context.Context, userID uint, notif model.Notification) error {
	wire := notificationWire{
		ID:        notif.ID,
		Type:      notif.Type,
		Title:     notif.Title,
		Body:      notif.Body,
		Data:      notif.Data,
		CreatedAt: notif.CreatedAt,
		ReadAt:    notif.ReadAt,
	}
	b, _ := json.Marshal(wire)
	n.hub.Send(userID, b)
	return nil
}
