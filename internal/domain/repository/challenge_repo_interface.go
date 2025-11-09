package repository

import "challenge-app/internal/domain/model"

type ChallengeRepository interface {
	// CRUD 
	CreateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error)
	GetChallengeByID(id uint) (*model.ChallengeModel, error)
	UpdateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error)
	DeleteChallenge(id uint) error
	StopChallenge(id uint) error

	// ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error)
	// ListPrivateChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	// ListInviteChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)

	// Challenges created by a specific user
	// ListByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error)

	// Challenges the user is participating in
	// ListByParticipant(userID uint, offset, limit int) ([]*model.ChallengeModel, error)

	// Challenges by category
	// ListByCategory(category string, offset, limit int) ([]*model.ChallengeModel, error)

	// Exists(challengeID uint) (bool, error)
	// IsCreator(userID, challengeID uint) (bool, error)
}
