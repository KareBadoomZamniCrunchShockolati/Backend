package repository

import "challenge-app/internal/domain/model"

type ChallengeRepository interface {
	// CRUD
	CreateChallenge(challenge *model.ChallengeModel) error
	GetChallengeByID(id uint) (*model.ChallengeModel, error)
	UpdateChallenge(challenge *model.ChallengeModel) error
	DeleteChallenge(id uint) error

	// List challenges in feed, respecting visibility rules
	ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error)
	ListPrivateChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	ListInviteChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)

	// List challenges created by a user
	ListByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error)

	// Challenge search / related
	ListByCategory(category string, offset, limit int) ([]*model.ChallengeModel, error)

	// Stop a challenge
	StopChallenge(id uint) error
}
