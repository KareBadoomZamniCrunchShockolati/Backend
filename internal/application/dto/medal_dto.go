package dto

import "time"

type MedalDTO struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	CategoryID  uint      `json:"category_id"`
	Description string    `json:"description"`
	AwardedAt   time.Time `json:"awarded_at,omitempty"`
}

type SelectMedalsDTO struct {
	Medals []string `json:"medals"` // for example "fitness:gold"
}
