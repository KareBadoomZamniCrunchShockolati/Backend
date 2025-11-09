package dto

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type CreateChallengeDTO struct {
	CreatorID       uint                      `json:"creator_id"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Category        string                   `json:"category"`
	MaxParticipants uint                     `json:"max_participants"`
	Visibility      enum.ChallengeVisibility `json:"visibility"`
	Rule            string                   `json:"rule"`
	StartTime       time.Time                `json:"start_time"`
	EndTime         time.Time                `json:"end_time"`
	Timezone        string                   `json:"timezone"`
	ImageURL        string                   `json:"image_url"`
}

type UpdateChallengeDTO struct {
	CreatorID       uint                      `json:"creator_id"`
	Title           *string                   `json:"title,omitempty"`
	Description     *string                   `json:"description,omitempty"`
	Category        *string                   `json:"category,omitempty"`
	MaxParticipants *uint                     `json:"max_participants,omitempty"`
	Rule            *string                   `json:"rule,omitempty"`
	Visibility      *enum.ChallengeVisibility `json:"visibility,omitempty"`
	EndTime         *time.Time                `json:"end_time,omitempty"`
	ImageURL        *string                   `json:"image_url,omitempty"`
	StartTime       *time.Time                `json:"start_time,omitempty"`
	Timezone        *string                   `json:"timezone,omitempty"`
}
