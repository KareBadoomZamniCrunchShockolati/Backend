package dto

import "time"

type NotificationResponseDTO struct {
	ID        uint           `json:"id"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
}

type NotificationListResponseDTO struct {
	Items      []NotificationResponseDTO `json:"items"`
	NextCursor *time.Time                `json:"next_cursor,omitempty"`
}
