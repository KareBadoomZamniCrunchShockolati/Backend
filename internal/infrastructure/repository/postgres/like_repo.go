package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) CreateLike(like *model.Like) error {
	likeEntity := &entity.LikeEntity{
		EntityType: string(like.EntityType),
		EntityID:   like.EntityID,
		UserID:     like.UserID,
	}

	result := r.db.Create(likeEntity)
	if result.Error != nil {
		// Check if it's a duplicate key error
		if strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "unique constraint") {
			return exception.NewConflictException("Like", "user_id", "USER_ALREADY_LIKED")
		}
		return exception.NewRepositoryError(result.Error)
	}

	return nil
}
func (r *LikeRepository) DeleteLike(entityType model.LikeType, entityID, userID uint) error {
	result := r.db.Where("entity_type = ? AND entity_id = ? AND user_id = ?",
		string(entityType), entityID, userID).Delete(&entity.LikeEntity{})

	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundException("Like",
			fmt.Sprintf("entity_type:%s,entity_id:%d,user_id:%d", entityType, entityID, userID),
			"LIKE_NOT_FOUND")
	}

	return nil
}

func (r *LikeRepository) IsUserLiked(entityType model.LikeType, entityID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.LikeEntity{}).
		Where("entity_type = ? AND entity_id = ? AND user_id = ?",
			entityType, entityID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LikeRepository) GetLikeCount(entityType model.LikeType, entityID uint) (uint, error) {
	var count int64
	err := r.db.Model(&entity.LikeEntity{}).
		Where("entity_type = ? AND entity_id = ?", string(entityType), entityID).
		Count(&count).Error

	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}
	return uint(count), nil
}
