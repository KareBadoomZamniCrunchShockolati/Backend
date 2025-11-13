package dto

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type CreateChallengeDTO struct {
	CreatorID       uint                     `json:"creator_id validate:"required"`
	Title           string                   `json:"title validate:"required,min=5,max=150"`
	Description     string                   `json:"description validate:"required,min=10,max=280"`
	CategoryID      uint                     `json:"category validate:"required"`
	MaxParticipants uint                     `json:"max_participants" validate:"omitempty,min=0"`
	Visibility      enum.ChallengeVisibility `json:"visibility validate:"required,oneof=1 2 3"`
	Rule            string                   `json:"rule validate:"required""`
	CommentsEnabled bool                     `json:"comments_enabled"`
	StartTime       time.Time                `json:"start_time"`
	EndTime         time.Time                `json:"end_time"`
	Timezone        string                   `json:"timezone"`
	ImageURL        string                   `json:"image_url"`
}

type UpdateChallengeDTO struct {
	CreatorID       *uint                     `json:"creator_id"`
	Title           *string                   `json:"title,omitempty"`
	Description     *string                   `json:"description,omitempty"`
	CategoryID      *uint                     `json:"category,omitempty"`
	CommentsEnabled *bool                     `json:"comments_enabled,omitempty"`
	IsStopped       *bool                     `json:"is_stopped,omitempty"`
	MaxParticipants *uint                     `json:"max_participants,omitempty"`
	Rule            *string                   `json:"rule,omitempty"`
	Visibility      *enum.ChallengeVisibility `json:"visibility,omitempty"`
	EndTime         *time.Time                `json:"end_time,omitempty"`
	ImageURL        *string                   `json:"image_url,omitempty"`
	StartTime       *time.Time                `json:"start_time,omitempty"`
	Timezone        *string                   `json:"timezone,omitempty"`
}

type ChallengeResponseDTO struct {
	ID                  uint                     `json:"id"`
	Title               string                   `json:"title"`
	Description         string                   `json:"description"`
	CategoryName        string                   `json:"category_name"`
	CreatorUsername     string                   `json:"creator_username"`
	MaxParticipants     uint                     `json:"max_participants"`
	CurrentParticipants int                      `json:"current_participants"`
	Visibility          enum.ChallengeVisibility `json:"visibility"`
	ImageURL            string                   `json:"image_url"`
	Rule                string                   `json:"recurrence_rule"`
	Timezone            string                   `json:"timezone"`
	StartTime           time.Time                `json:"start_time"`
	EndTime             *time.Time               `json:"end_time,omitempty"`
	IsStopped           bool                     `json:"is_stopped"`
	CommentsEnabled     bool                     `json:"comments_enabled"`
	CreatedAt           time.Time                `json:"created_at"`
	IsUserParticipating bool                     `json:"is_user_participating"`
}
