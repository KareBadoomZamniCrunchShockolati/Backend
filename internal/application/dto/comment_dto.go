package dto

import (
	"time"
)

type CommentRequestDTO struct {
	EntityType string `json:"entity_type" validate:"required,oneof=challenge post"`
	EntityID   uint   `json:"entity_id" validate:"required"`
	Content    string `json:"content" validate:"required,min=1,max=1000"`
	ParentID   *uint  `json:"parent_id"`
}

type CommentResponseDTO struct {
	ID         uint                  `json:"id"`
	EntityType string                `json:"entity_type"`
	EntityID   uint                  `json:"entity_id"`
	UserID     uint                  `json:"user_id"`
	Username   string                `json:"username"`
	Content    string                `json:"content"`
	ParentID   *uint                 `json:"parent_id,omitempty"`
	LikeCount  uint                  `json:"like_count"`
	IsLiked    bool                  `json:"is_liked"`
	CreatedAt  time.Time             `json:"created_at"`
	Replies    []*CommentResponseDTO `json:"replies,omitempty"`
}
