package repository

import "challenge-app/internal/domain/model"

type ChallengeRepository interface {
	// CRUD
	CreateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error)
	GetChallengeByID(id uint) (*model.ChallengeModel, error)
	GetAllChallenges() ([]*model.ChallengeModel, error)
	UpdateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error)
	DeleteChallenge(id uint) error
	StopChallenge(id uint) error

	// List challenges with filters
	ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error)
	ListChallengesByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	ListChallengesByCategory(categoryID uint, offset, limit int) ([]*model.ChallengeModel, error)
	SearchChallengesByCategory(categoryName string, offset, limit int) ([]*model.ChallengeModel, error)
	ListChallengesByParticipant(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	GetMutualFollowersInChallenge(userID, challengeID uint) ([]*model.UserModel, error)

	// Check if user is challenge creator
	IsChallengeCreator(challengeID, userID uint) (bool, error)
}
