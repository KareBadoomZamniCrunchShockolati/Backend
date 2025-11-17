package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"log"

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
	return toChallengeModelWithID(*entity, challenge), nil
}

func (r *ChallengeRepository) GetChallengeByID(id uint) (*model.ChallengeModel, error) {
	var entity entity.ChallengeEntity
	result := r.db.Preload("Participants").Preload("Comments").First(&entity, id)
	if result.Error != nil {
		log.Println("here in repo")
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
		Where("cp.user_id = ? AND cp.status IN ?", userID, []uint{uint(enum.StatusJoined), uint(enum.StatusPending)}).
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
