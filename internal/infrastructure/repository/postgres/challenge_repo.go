package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"fmt"

	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{db: db}
}

func toChallengeEntity(m *model.ChallengeModel) *entity.ChallengeEntity {
	return &entity.ChallengeEntity{
		Title:           m.Title,
		Description:     m.Description,
		Category:        m.Category,
		CreatorID:       m.CreatorID,
		MaxParticipants: m.MaxParticipants,
		Visibility:      uint(m.Visibility),
		Rule:            m.Rule,
		StartTime:       &m.StartTime,
		EndTime:         m.EndTime,
		Timezone:        m.Timezone,
		ImageURL:        m.ImageURL,
		Stopped:         m.Stopped,
	}
}

func toChallengeModel(e *entity.ChallengeEntity) *model.ChallengeModel {
	return &model.ChallengeModel{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		Category:        e.Category,
		CreatorID:       e.CreatorID,
		MaxParticipants: e.MaxParticipants,
		Visibility:      enum.ChallengeVisibility(e.Visibility),
		Rule:            e.Rule,
		EndTime:         e.EndTime,
		Timezone:        e.Timezone,
		ImageURL:        e.ImageURL,
		Stopped:         e.Stopped,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		// Participants:  toParticipantModels(e.Participants),
		// Comments:      toCommentModels(e.Comments),
	}
}

func toParticipantModels(entities []entity.ChallengeParticipantEntity) []*model.ChallengeParticipant {
	var participants []*model.ChallengeParticipant
	for _, e := range entities {
		participants = append(participants, &model.ChallengeParticipant{
			ID:          e.ID,
			ChallengeID: e.ChallengeID,
			UserID:      e.UserID,
			Status:      e.Status,
			CreatedAt:   e.CreatedAt,
		})
	}
	return participants
}

func toCommentModels(entities []entity.ChallengeCommentEntity) []*model.ChallengeComment {
	var comments []*model.ChallengeComment
	for _, e := range entities {
		comments = append(comments, &model.ChallengeComment{
			ID:          e.ID,
			ChallengeID: e.ChallengeID,
			UserID:      e.UserID,
			Content:     e.Content,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		})
	}
	return comments
}
func (r *ChallengeRepository) CreateChallenge(ch *model.ChallengeModel) (*model.ChallengeModel, error) {
	chEntity := toChallengeEntity(ch)
	if err := r.db.Create(chEntity).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	ch.ID = chEntity.ID
	return ch, nil
}

func (r *ChallengeRepository) GetChallengeByID(id uint) (*model.ChallengeModel, error) {
	var chEntity entity.ChallengeEntity
	err := r.db.Preload("Participants").Preload("Comments").First(&chEntity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
		}
		return nil, exception.NewRepositoryError(err)
	}
	return toChallengeModel(&chEntity), nil
}

func (r *ChallengeRepository) GetAllChallenges() ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	err := r.db.Find(&chEntities).Error
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	var challenges []*model.ChallengeModel
	for _, e := range chEntities {
		challenges = append(challenges, toChallengeModel(&e))
	}
	return challenges, nil
}

func (r *ChallengeRepository) UpdateChallenge(ch *model.ChallengeModel) (*model.ChallengeModel, error) {
	var chEntity entity.ChallengeEntity
	if err := r.db.First(&chEntity, ch.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", ch.ID), "CHALLENGE_NOT_FOUND")
		}
		return nil, exception.NewRepositoryError(err)
	}

	// Map changes from model to entity
	chEntity = *toChallengeEntity(ch)

	if err := r.db.Save(&chEntity).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	return toChallengeModel(&chEntity), nil
}

func (r *ChallengeRepository) DeleteChallenge(id uint) error {
	if err := r.db.Delete(&entity.ChallengeEntity{}, id).Error; err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (r *ChallengeRepository) StopChallenge(id uint) error {
	if err := r.db.Model(&entity.ChallengeEntity{}).Where("id = ?", id).Update("stopped", true).Error; err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}
