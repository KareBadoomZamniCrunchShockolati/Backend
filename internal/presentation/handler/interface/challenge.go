package handler

import (
	"github.com/gin-gonic/gin"
)

type ChallengeHandler interface {
	CreateChallenge(c *gin.Context)
	GetChallengeByID(c *gin.Context)
	UpdateChallenge(c *gin.Context)
	DeleteChallenge(c *gin.Context)

	// Listing challenges
	// ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error)
	// ListPrivateChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	// ListInviteChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	// ListByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	// ListByCategory(category string, offset, limit int) ([]*model.ChallengeModel, error)
}
