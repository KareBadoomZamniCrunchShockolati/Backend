package repository

import "challenge-app/internal/domain/model"

type LikeRepository interface {
	CreateLike(like *model.ChallengeLike) error
	DeleteLike(challengeID, userID uint) error
	IsUserLikedChallenge(challengeID, userID uint) (bool, error)
	GetLikeCount(challengeID uint) (uint, error)
}