package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"log"
	"fmt"
	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{db: db}
}

func (r *ChallengeRepository) CreateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error) {
	entity := toChallengeEntity(challenge)
	result := r.db.Create(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	log.Println("Created Challenge with StartTime:", challenge.StartTime)
	return toChallengeModelWithID(*entity, challenge), nil
}

func (r *ChallengeRepository) GetChallengeByID(id uint) (*model.ChallengeModel, error) {
	var entity entity.ChallengeEntity
	result := r.db.Preload("Participants").Preload("Comments").First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(result.Error)
	}
	return toChallengeModel(&entity), nil
}

func (r *ChallengeRepository) GetAllChallenges() ([]*model.ChallengeModel, error) {
	var entities []entity.ChallengeEntity
	result := r.db.Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	challenges := make([]*model.ChallengeModel, len(entities))
	for i, e := range entities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) UpdateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error) {
	entity := toChallengeEntity(challenge)
	entity.ID = challenge.ID
	log.Println("Updating Challenge with StartTime:", challenge.StartTime)
	result := r.db.Save(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	challenge.UpdatedAt = entity.UpdatedAt
	return challenge, nil
}

func (r *ChallengeRepository) DeleteChallenge(id uint) error {
	result := r.db.Delete(&entity.ChallengeEntity{}, id)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

func (r *ChallengeRepository) StopChallenge(id uint) error {
	result := r.db.Model(&entity.ChallengeEntity{}).Where("id = ?", id).Update("is_stopped", true)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

// List methods
func (r *ChallengeRepository) ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	result := r.db.Where("visibility = ? AND is_stopped = ?", uint(enum.VisibilityPublic), false).
		Offset(offset).Limit(limit).Find(&chEntities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	challenges := make([]*model.ChallengeModel, len(chEntities))
	for i, e := range chEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) ListChallengesByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	result := r.db.Where("creator_id = ?", userID).
		Offset(offset).Limit(limit).Find(&chEntities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	challenges := make([]*model.ChallengeModel, len(chEntities))
	for i, e := range chEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) ListChallengesByCategory(categoryID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	result := r.db.Where("category_id = ?", categoryID).
		Offset(offset).Limit(limit).Find(&chEntities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	challenges := make([]*model.ChallengeModel, len(chEntities))
	for i, e := range chEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) ListChallengesByParticipant(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	result := r.db.Joins("JOIN challenge_participants cp ON cp.challenge_id = challenges.id").
		Where("cp.user_id = ? AND cp.status IN ?", userID, []uint{uint(enum.StatusJoined)}).
		Offset(offset).Limit(limit).Find(&chEntities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	challenges := make([]*model.ChallengeModel, len(chEntities))
	for i, e := range chEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) SearchChallengesByCategory(categoryName string, offset, limit int) ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	result := r.db.Joins("JOIN challenge_categories cc ON cc.id = challenges.category_id").
		Where("cc.name ILIKE ?", "%"+categoryName+"%").
		Offset(offset).Limit(limit).Find(&chEntities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	challenges := make([]*model.ChallengeModel, len(chEntities))
	for i, e := range chEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) IsChallengeCreator(challengeID, userID uint) (bool, error) {
	var chEntity entity.ChallengeEntity
	result := r.db.Select("creator_id").Where("id = ?", challengeID).First(&chEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, exception.NewRepositoryError(result.Error)
	}
	return chEntity.CreatorID == userID, nil
}

func (r *ChallengeRepository) GetMutualFollowersInChallenge(userID, challengeID uint) ([]*model.UserModel, error) {
	var userEntities []entity.UserEntity

	err := r.db.Table("users u").
		Joins("INNER JOIN challenge_participants cp ON u.id = cp.user_id").
		Joins("INNER JOIN follows f ON u.id = f.following_id").
		Where("cp.challenge_id = ? AND f.follower_id = ? AND cp.status IN ? AND f.status = ?",
			challengeID,
			userID,
			[]uint{uint(enum.StatusJoined), uint(enum.StatusPending)},
			"active").
		Select("u.id, u.username, u.email, u.bio, u.verified").
		Find(&userEntities).Error

	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	userModels := make([]*model.UserModel, len(userEntities))
	for i, u := range userEntities {
		userModels[i] = &model.UserModel{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Bio:      u.Bio,
			Verified: u.Verified,
		}
	}

	return userModels, nil
}

func (r *ChallengeRepository) CreateLike(like *model.ChallengeLike) error {
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

func (r *ChallengeRepository) DeleteLike(challengeID, userID uint) error {
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

func (r *ChallengeRepository) IsUserLikedChallenge(challengeID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.ChallengeLikeEntity{}).
		Where("challenge_id = ? AND user_id = ?", challengeID, userID).
		Count(&count).Error
	if err != nil {
		return false, exception.NewRepositoryError(err)
	}
	return count > 0, nil
}

func (r *ChallengeRepository) GetLikeCount(challengeID uint) (uint, error) {
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

// Helpers
func toChallengeEntity(m *model.ChallengeModel) *entity.ChallengeEntity {
	return &entity.ChallengeEntity{
		Title:           m.Title,
		Description:     m.Description,
		CategoryID:      m.CategoryID,
		CreatorID:       m.CreatorID,
		MaxParticipants: m.MaxParticipants,
		Visibility:      uint(m.Visibility),
		Rule:            m.Rule,
		StartTime:       &m.StartTime,
		EndTime:         m.EndTime,
		Timezone:        m.Timezone,
		ImageURL:        m.ImageURL,
		IsStopped:       m.IsStopped,
		CommentsEnabled: m.CommentsEnabled,
	}
}

func toChallengeModel(e *entity.ChallengeEntity) *model.ChallengeModel {
	return &model.ChallengeModel{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		CategoryID:      e.CategoryID,
		CreatorID:       e.CreatorID,
		MaxParticipants: e.MaxParticipants,
		Visibility:      enum.ChallengeVisibility(e.Visibility),
		Rule:            e.Rule,
		EndTime:         e.EndTime,
		StartTime:       *e.StartTime,
		Timezone:        e.Timezone,
		ImageURL:        e.ImageURL,
		IsStopped:       e.IsStopped,
		CommentsEnabled: e.CommentsEnabled,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func toChallengeModelWithID(e entity.ChallengeEntity, original *model.ChallengeModel) *model.ChallengeModel {
	original.ID = e.ID
	original.CreatedAt = e.CreatedAt
	original.UpdatedAt = e.UpdatedAt
	return original
}
