package ws

import (
	"context"
	"encoding/json"
	"log"

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
	TitleKey  string                 `json:"title_key"`
	BodyKey   string                 `json:"body_key"`
	Data      map[string]any         `json:"data,omitempty"`
	CreatedAt string                 `json:"created_at"` 
	ReadAt    *string                `json:"read_at,omitempty"`
}

func (n *WSNotifierImpl) Push(ctx context.Context, userID uint, notif model.Notification) error {
	created := notif.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	var readAt *string
	if notif.ReadAt != nil {
		s := notif.ReadAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		readAt = &s
	}

	log.Println("PUSH to user:", userID, "notif id:", notif.ID)

	wire := notificationWire{
		ID:        notif.ID,
		Type:      notif.Type,
		TitleKey:  notif.TitleKey,
		BodyKey:   notif.BodyKey,
		Data:      notif.Data,
		CreatedAt: created,
		ReadAt:    readAt,
	}

	b, err := json.Marshal(wire)
	if err != nil {
		return err
	}

	n.hub.Send(userID, b)
	return nil
}
