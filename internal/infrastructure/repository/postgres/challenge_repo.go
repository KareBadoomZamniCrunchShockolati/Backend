package postgres

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"challenge-app/internal/domain/exception"
	"gorm.io/gorm"
	"fmt"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepo(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{db: db}
}


//Helper functions
func toParticipantEntities(ms []*model.ChallengeParticipant) []*entity.ChallengeParticipantEntity {
	entities := make([]*entity.ChallengeParticipantEntity, len(ms))
	for i, m := range ms {
		entities[i] = toParticipantEntity(m)
	}
	return entities
}

func toCommentEntities(ms []*model.ChallengeComment) []*entity.ChallengeCommentEntity {
	entities := make([]*entity.ChallengeCommentEntity, len(ms))
	for i, m := range ms {
		entities[i] = toCommentEntity(m)
	}
	return entities
}

func toChallengeEntity(ch *model.ChallengeModel) *entity.ChallengeEntity {
	return &entity.ChallengeEntity{
		Title:           ch.Title,
		Description:     ch.Description,
		Category:        ch.Category,
		CreatorID:       ch.CreatorID,
		MaxParticipants: ch.MaxParticipants,
		Visibility:      ch.Visibility,
		ImageURL:        ch.ImageURL,
		Rule:            ch.Rule,
		Timezone:        ch.Timezone,
		StartTime:       &ch.StartTime,
		EndTime:         ch.EndTime,
		Stopped:         ch.Stopped,
		// Relations
		Participants: toParticipantEntities(ch.Participants),
		Comments:     toCommentEntities(ch.Comments),
	}
}

func toParticipantModels(es []*entity.ChallengeParticipantEntity) []*model.ChallengeParticipant {
	models := make([]*model.ChallengeParticipant, 0, len(es))
	for i := range es {
		models = append(models, toParticipantModel(es[i]))
	}
	return models
}

func toCommentModels(es []*entity.ChallengeCommentEntity) []*model.ChallengeComment {
	models := make([]*model.ChallengeComment, 0, len(es))
	for i := range es {
		models = append(models, toCommentModel(es[i]))
	}
	return models
}

func toChallengeModel(e *entity.ChallengeEntity) *model.ChallengeModel {
	ch := &model.ChallengeModel{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		Category:        e.Category,
		CreatorID:       e.CreatorID,
		MaxParticipants: e.MaxParticipants,
		Visibility:      e.Visibility,
		ImageURL:        e.ImageURL,
		Rule:            e.Rule,
		Timezone:        e.Timezone,
		Stopped:         e.Stopped,
		StartTime:       *e.StartTime,
		EndTime:         e.EndTime,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		Participants:    toParticipantModels(e.Participants),
		Comments:        toCommentModels(e.Comments),
	}
	return ch
}

func toChallengeModels(es []entity.ChallengeEntity) []*model.ChallengeModel {
	models := make([]*model.ChallengeModel, 0, len(es))
	for i := range es {
		models = append(models, toChallengeModel(&es[i]))
	}
	return models
}




func (r *ChallengeRepository) CreateChallenge(ch *model.ChallengeModel) (*model.ChallengeModel, error) {
	chEntity := toChallengeEntity(ch)

	if err := r.db.Create(chEntity).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("failed to create challenge"), err)
	}

	ch.ID = chEntity.ID
	return ch, nil
}

func (r *ChallengeRepository) GetChallengeByID(id uint) (*model.ChallengeModel, error) {
	var chEntity entity.ChallengeEntity
	err := r.db.Preload("Participants").Preload("Comments").First(&chEntity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE NOT FOUND")
		}
		return nil, exception.NewRepositoryError(fmt.Sprintf("GetChallengeByID %d", id), err)
	}

	return toChallengeModel(&chEntity), nil
}

func (r *ChallengeRepository) UpdateChallenge(ch *model.ChallengeModel) (*model.ChallengeModel, error) {
	var chEntity entity.ChallengeEntity
	if err := r.db.First(&chEntity, ch.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("Challenge", chEntity.Title, "CHALLENGE NOT FOUND")
		}
		return nil, exception.NewRepositoryError("failed to update challenge", err)
	}

	chEntity = *toChallengeEntity(ch) 
	if err := r.db.Save(&chEntity).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("UpdateChallenge %d", ch.ID), err)
	}
	return toChallengeModel(&chEntity), nil
}

func (r *ChallengeRepository) DeleteChallenge(id uint) error {
	if err := r.db.Delete(&entity.ChallengeEntity{}, id).Error; err != nil {
		return exception.NewRepositoryError(fmt.Sprintf("DeleteChallenge %d", id), err)
	}
	return nil
}

func (r *ChallengeRepository) ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error) {
	var entities []entity.ChallengeEntity
	if err := r.db.Where("visibility = ?", "public").Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("ListPublicChallenges"), err)
	}
	return toChallengeModels(entities), nil
}

func (r *ChallengeRepository) ListPrivateChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var entities []entity.ChallengeEntity
	if err := r.db.Where("visibility = ? AND creator_id = ?", "private", userID).Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("ListPrivateChallengesForUser %d", userID), err)
	}
	return toChallengeModels(entities), nil
}

func (r *ChallengeRepository) ListInviteChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var entities []entity.ChallengeEntity
	if err := r.db.Table("challenge_participants").Select("challenges.*").
		Joins("JOIN challenges ON challenges.id = challenge_participants.challenge_id").
		Where("challenge_participants.user_id = ?", userID).
		Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("ListInviteChallengesForUser %d", userID), err)
	}
	return toChallengeModels(entities), nil
}

func (r *ChallengeRepository) ListByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var entities []entity.ChallengeEntity
	if err := r.db.Where("creator_id = ?", userID).Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("ListByCreator %d", userID), err)
	}

	return toChallengeModels(entities), nil
}

func (r *ChallengeRepository) ListByCategory(category string, offset, limit int) ([]*model.ChallengeModel, error) {
	var entities []entity.ChallengeEntity
	if err := r.db.Where("category = ?", category).Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("ListByCategory %s", category), err)
	}
	return toChallengeModels(entities), nil
}

func (r *ChallengeRepository) StopChallenge(id uint) error {
	if err := r.db.Model(&entity.ChallengeEntity{}).Where("id = ?", id).Update("stopped", true).Error; err != nil {
		return exception.NewRepositoryError(fmt.Sprintf("StopChallenge %d", id), err)
	}
	return nil
}
