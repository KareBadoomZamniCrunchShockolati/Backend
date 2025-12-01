package dto

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type CreateChallengeDTO struct {
	CreatorID       uint                     `json:"creator_id" validate:"required"`
	Title           string                   `json:"title" validate:"required,min=5,max=150"`
	Description     string                   `json:"description" validate:"required,min=10,max=280"`
	CategoryID      uint                     `json:"category_id" validate:"required"`
	MaxParticipants uint                     `json:"max_participants" validate:"omitempty,min=0"`
	Visibility      enum.ChallengeVisibility `json:"visibility" validate:"required,oneof=public private invite"`
	Location        string                   `json:"location" validate:"required"`
	Goal            int                      `json:"goal" validate:"required,default=1"`
	Rule            string                   `json:"rule" validate:"required"`
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
	Location        *string                   `json:"location,omitempty"`
	Goal            *int                      `json:"goal,omitempty"`
	EndTime         *time.Time                `json:"end_time,omitempty"`
	ImageURL        *string                   `json:"image_url,omitempty"`
	StartTime       *time.Time                `json:"start_time,omitempty"`
	Timezone        *string                   `json:"timezone,omitempty"`
}

type ChallengePreviewDTO struct {
	ID                  uint                     `json:"id"`
	Title               string                   `json:"title"`
	Description         string                   `json:"description"`
	Rule                string                   `json:"recurrence_rule"`
	CategoryName        string                   `json:"category_name"`
	CreatorUsername     string                   `json:"creator_username"`
	CreatorID           uint                     `json:"creator_id"`
	Visibility          enum.ChallengeVisibility `json:"visibility"`
	Location            string                   `json:"location"`
	Goal                int                      `json:"goal"`
	ImageURL            string                   `json:"image_url"`
	MaxParticipants     uint                     `json:"max_participants"`
	CurrentParticipants int                      `json:"current_participants"`
	LikeCount           uint                     `json:"like_count"`
	CommentCount        uint                     `json:"comment_count"`
	StartTime           time.Time                `json:"start_time"`
	EndTime             *time.Time               `json:"end_time,omitempty"`
	Timezone            string                   `json:"timezone"`
	CreatedAt           time.Time                `json:"created_at"`
	IsUserParticipating bool                     `json:"is_user_participating"`
	IsUserLiked         bool                     `json:"is_user_liked"`
	MutualParticipants  []UserPreviewDTO         `json:"mutual_participants"`
}

type ChallengeDetailDTO struct {
	ChallengePreviewDTO
	CommentsEnabled bool                     `json:"comments_enabled"`
	Participants    []ParticipantResponseDTO `json:"participants,omitempty"`
	Comments        []CommentResponseDTO     `json:"comments,omitempty"`
}
