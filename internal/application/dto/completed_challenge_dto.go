package dto

import "time"

type CompletedChallengeDTO struct {
	ChallengePreviewDTO
	CompletedAt     time.Time `json:"completed_at"`
	CompletionRate  float64   `json:"completion_rate"`
	ChallengeGoal   int       `json:"challenge_goal"`
	AverageProgress float64   `json:"average_progress"`
}

type UserCompletionStatsDTO struct {
	TotalChallenges     int     `json:"total_challenges"`
	CompletedChallenges int     `json:"completed_challenges"`
	CompletionRate      float64 `json:"completion_rate"`
}
