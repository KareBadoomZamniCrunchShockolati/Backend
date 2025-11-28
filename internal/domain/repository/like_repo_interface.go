package repository

import "challenge-app/internal/domain/model"

type LikeRepository interface {
	CreateLike(like *model.Like) error
	DeleteLike(entityType model.LikeType, entityID, userID uint) error
	IsUserLiked(entityType model.LikeType, entityID, userID uint) (bool, error)
	GetLikeCount(entityType model.LikeType, entityID uint) (uint, error)
}
