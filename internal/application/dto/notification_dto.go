package dto

import "time"

type NotificationResponseDTO struct {
	ID        uint           `json:"id"`
	Type      string         `json:"type"`
	TitleKey  string         `json:"title_key"`
	BodyKey   string         `json:"body_key"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
}

type NotificationListResponseDTO struct {
	Items      []NotificationResponseDTO `json:"items"`
	NextCursor *time.Time                `json:"next_cursor,omitempty"`
}
