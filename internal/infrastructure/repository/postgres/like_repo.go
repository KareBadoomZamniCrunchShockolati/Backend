package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"fmt"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) CreateLike(like *model.ChallengeLike) error {
	likeEntity := &entity.ChallengeLikeEntity{
		ChallengeID: like.ChallengeID,
		UserID:      like.UserID,
	}
	result := r.db.Create(likeEntity)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	result = r.db.Model(&entity.ChallengeEntity{}).
		Where("id = ?", like.ChallengeID).
		Update("like_count", gorm.Expr("like_count + ?", 1))
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}

	return nil
}

func (r *LikeRepository) DeleteLike(challengeID, userID uint) error {
	result := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).Delete(&entity.ChallengeLikeEntity{})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundException("Like", fmt.Sprintf("challenge_id:%d,user_id:%d", challengeID, userID), "LIKE_NOT_FOUND")
	}

	result = r.db.Model(&entity.ChallengeEntity{}).
		Where("id = ?", challengeID).
		Update("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", 1))
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}

	return nil
}

func (r *LikeRepository) IsUserLikedChallenge(challengeID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.ChallengeLikeEntity{}).
		Where("challenge_id = ? AND user_id = ?", challengeID, userID).
		Count(&count).Error
	if err != nil {
		return false, exception.NewRepositoryError(err)
	}
	return count > 0, nil
}

func (r *LikeRepository) GetLikeCount(challengeID uint) (uint, error) {
	var challenge entity.ChallengeEntity
	err := r.db.Select("like_count").Where("id = ?", challengeID).First(&challenge).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, exception.NewRepositoryError(err)
	}
	return challenge.LikeCount, nil
}
