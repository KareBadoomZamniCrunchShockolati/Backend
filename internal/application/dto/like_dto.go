package dto

import "time"

type LikeRequestDTO struct {
	EntityType string `json:"entity_type" validate:"required,oneof=challenge post comment"`
	EntityID   uint   `json:"entity_id" validate:"required"`
}

type LikeResponseDTO struct {
	ID         uint      `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   uint      `json:"entity_id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

type LikeStatsDTO struct {
	LikeCount uint `json:"like_count"`
	IsLiked   bool `json:"is_liked"`
}
