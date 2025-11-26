package repository

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/model"
)

type ChallengeRepository interface {
	// CRUD
	CreateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error)
	GetChallengeByID(id uint, userID uint) (*model.ChallengeModel, error)
	GetAllChallenges(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	UpdateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error)
	DeleteChallenge(id uint) error
	StopChallenge(id uint) error

	// List challenges with filters
	ListDiscoverableChallenges(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListPublicChallenges(offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCreatorID(creatorID, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCreatorUsername(username string, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCategoryID(categoryID, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCategoryName(name string, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByParticipant(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	ListChallengesByParticipantCount(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByLikeCount(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesStartingSoon(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListTopCreatorsChallenge(offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesJoinedByUser(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	SearchChallenges(query string, visibility []enum.ChallengeVisibility, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)

	GetMutualFollowersInChallenge(userID, challengeID uint) ([]*model.UserModel, error)

	// Like
	CreateLike(like *model.ChallengeLike) error
	DeleteLike(challengeID, userID uint) error
	IsUserLikedChallenge(challengeID, userID uint) (bool, error)
	GetLikeCount(challengeID uint) (uint, error)

	// Check if user is challenge creator
	IsChallengeCreator(challengeID, userID uint) (bool, error)
}
