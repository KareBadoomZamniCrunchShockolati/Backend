package serviceinterface

import (
	"challenge-app/internal/application/dto"
)

type ChallengeCompletionServicer interface {
	// Process ended challenges and mark completions
	ProcessEndedChallenges() (int, error)

	// Check if a user has completed a specific challenge
	CheckUserCompletion(userID, challengeID uint) (bool, float64, error)

	// Get user's completed challenges
	GetUserCompletedChallenges(userID uint, offset, limit int) ([]*dto.CompletedChallengeDTO, error)

	// Get user completion statistics
	GetUserCompletionStats(userID uint) (*dto.UserCompletionStatsDTO, error)

	// Get challenge completion rate
	GetChallengeCompletionRate(challengeID uint) (float64, error)
}
