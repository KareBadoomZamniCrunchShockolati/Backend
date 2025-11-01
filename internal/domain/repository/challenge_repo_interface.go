package repository

import "challenge-app/internal/domain/model"

type ChallengeRepository interface {
	// CRUD
	Create(challenge *model.Challenge) error
	GetByID(id uint) (*model.Challenge, error)
	Update(challenge *model.Challenge) error
	Delete(id uint) error

	// List challenges in feed, respecting visibility rules
	ListPublicChallenges(offset, limit int) ([]*model.Challenge, error)
	ListPrivateChallengesForUser(userID uint, offset, limit int) ([]*model.Challenge, error)
	ListInviteChallengesForUser(userID uint, offset, limit int) ([]*model.Challenge, error)

	// List challenges created by a user
	ListByCreator(userID uint, offset, limit int) ([]*model.Challenge, error)

	// Challenge search / related
	ListByCategory(category string, offset, limit int) ([]*model.Challenge, error)

	// Stop a challenge
	StopChallenge(id uint) error
}
