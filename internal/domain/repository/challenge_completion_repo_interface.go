package repository

import "challenge-app/internal/domain/model"

type ChallengeCompletionRepository interface {
	// CRUD operations
	CreateCompletion(completion *model.ChallengeCompletion) (*model.ChallengeCompletion, error)
	GetCompletion(userID, challengeID uint) (*model.ChallengeCompletion, error)
	UpdateCompletion(completion *model.ChallengeCompletion) (*model.ChallengeCompletion, error)
	DeleteCompletion(userID, challengeID uint) error

	// Batch operations
	MarkChallengeCompleted(userID, challengeID uint) error
	GetCompletedChallengesByUser(userID uint, offset, limit int) ([]*model.ChallengeCompletion, error)
	GetChallengesToProcess() ([]*model.ChallengeModel, error)

	// Stats
	GetCompletionRate(challengeID uint) (float64, error) // Percentage of participants who completed
	GetUserCompletionStats(userID uint) (total, completed int, err error)
}
