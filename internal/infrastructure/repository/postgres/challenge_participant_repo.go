package postgres

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"challenge-app/internal/domain/exception"
	"gorm.io/gorm"
	"fmt"
)

type ChallengeParticipantRepository struct {
	db *gorm.DB
}

func NewChallengeParticipantRepo(db *gorm.DB) *ChallengeParticipantRepository {
	return &ChallengeParticipantRepository{db: db}
}

//Helper functions
func toParticipantEntity(p *model.ChallengeParticipant) *entity.ChallengeParticipantEntity {
	return &entity.ChallengeParticipantEntity{
		Model:       gorm.Model{ID: p.ID, CreatedAt: p.CreatedAt},
		ChallengeID: p.ChallengeID,
		UserID:      p.UserID,
		Status:      p.Status,
	}
}

func toParticipantModel(e *entity.ChallengeParticipantEntity) *model.ChallengeParticipant {
	return &model.ChallengeParticipant{
		ID:          e.ID,
		ChallengeID: e.ChallengeID,
		UserID:      e.UserID,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
	}
}


// Add a participant
func (r *ChallengeParticipantRepository) AddParticipant(p *model.ChallengeParticipant) error {
	pEntity := toParticipantEntity(p)
	if err := r.db.Create(pEntity).Error; err != nil {
		return exception.NewRepositoryError("failed to add participant", err)
	}
	p.ID = pEntity.ID
	return nil
}

func (r *ChallengeParticipantRepository) UpdateParticipantStatus(p *model.ChallengeParticipant) error {
	if err := r.db.Model(&entity.ChallengeParticipantEntity{}).
		Where("challenge_id = ? AND user_id = ?", p.ChallengeID, p.UserID).
		Update("status", p.Status.String()).Error; err != nil {
		return exception.NewRepositoryError("failed to update participant status", err)
	}
	return nil
}

func (r *ChallengeParticipantRepository) RemoveParticipant(challengeID, userID uint) error {
	if err := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).
		Delete(&entity.ChallengeParticipantEntity{}).Error; err != nil {
		return exception.NewRepositoryError("failed to remove participant", err)
	}
	return nil
}

func (r *ChallengeParticipantRepository) GetParticipant(challengeID, userID uint) (*model.ChallengeParticipant, error) {
	var pEntity entity.ChallengeParticipantEntity
	err := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).First(&pEntity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("ChallengeParticipant", fmt.Sprintf("%d-%d", challengeID, userID), "PARTICIPANT NOT FOUND")
		}
		return nil, exception.NewRepositoryError("failed to get participant", err)
	}
	return toParticipantModel(&pEntity), nil
}

func (r *ChallengeParticipantRepository) ListParticipants(challengeID uint) ([]*model.ChallengeParticipant, error) {
	var entities []entity.ChallengeParticipantEntity
	if err := r.db.Where("challenge_id = ?", challengeID).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError("failed to list participants", err)
	}

	var participants []*model.ChallengeParticipant
	for _, e := range entities {
		participants = append(participants, toParticipantModel(&e))
	}
	return participants, nil
}

func (r *ChallengeParticipantRepository) ListParticipantsInUserFollowings(challengeID, userID uint) ([]*model.ChallengeParticipant, error) {
	var entities []entity.ChallengeParticipantEntity
	if err := r.db.Where("challenge_id = ? AND user_id IN (SELECT followed_user_id FROM user_followings WHERE follower_user_id = ?)", challengeID, userID).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError("failed to list participants in user followings", err)
	}

	var participants []*model.ChallengeParticipant
	for _, e := range entities {
		participants = append(participants, toParticipantModel(&e))
	}
	return participants, nil
}

func (r *ChallengeParticipantRepository) ListUserChallenges(userID uint) ([]*model.ChallengeParticipant, error) {
	var entities []entity.ChallengeParticipantEntity
	if err := r.db.Where("user_id = ?", userID).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError("failed to list user challenges", err)
	}

	var participants []*model.ChallengeParticipant
	for _, e := range entities {
		participants = append(participants, toParticipantModel(&e))
	}
	return participants, nil
}