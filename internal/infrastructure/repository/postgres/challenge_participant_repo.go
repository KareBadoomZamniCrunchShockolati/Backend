package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type ChallengeParticipantRepository struct {
	db *gorm.DB
}

func NewChallengeParticipantRepository(db *gorm.DB) *ChallengeParticipantRepository {
	return &ChallengeParticipantRepository{db: db}
}

func (r *ChallengeParticipantRepository) CreateParticipant(participant *model.ChallengeParticipant) (*model.ChallengeParticipant, error) {
	entity := toChallengeParticipantEntity(participant)
	result := r.db.Create(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	return toChallengeParticipantModelWithID(*entity, participant), nil
}

func (r *ChallengeParticipantRepository) GetParticipant(challengeID, userID uint) (*model.ChallengeParticipant, error) {
	var entity entity.ChallengeParticipantEntity
	result := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).First(&entity)
	return handleGetEntityWithCondition(result.Error, func() *model.ChallengeParticipant {
		return toChallengeParticipantModel(&entity)
	})
}

func (r *ChallengeParticipantRepository) GetParticipantsByChallenge(challengeID uint) ([]*model.ChallengeParticipant, error) {
	var entities []entity.ChallengeParticipantEntity
	result := r.db.Where("challenge_id = ?", challengeID).Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	participants := make([]*model.ChallengeParticipant, len(entities))
	for i, e := range entities {
		participants[i] = toChallengeParticipantModel(&e)
	}
	return participants, nil
}

func (r *ChallengeParticipantRepository) UpdateParticipant(participant *model.ChallengeParticipant) (*model.ChallengeParticipant, error) {
	entity := toChallengeParticipantEntity(participant)
	entity.ID = participant.ID

	result := r.db.Save(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	participant.UpdatedAt = entity.UpdatedAt
	return participant, nil
}

func (r *ChallengeParticipantRepository) DeleteParticipant(challengeID, userID uint) error {
	result := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).Delete(&entity.ChallengeParticipantEntity{})
	return handleDeleteResult(result.Error, "Participant", 0)
}

func (r *ChallengeParticipantRepository) GetParticipantCount(challengeID uint) (int, error) {
	var count int64
	err := r.db.Model(&entity.ChallengeParticipantEntity{}).
		Where("challenge_id = ? AND status IN ?", challengeID, []uint{uint(enum.StatusJoined), uint(enum.StatusPending)}).
		Count(&count).Error
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}
	return int(count), nil
}

func (r *ChallengeParticipantRepository) IsUserParticipant(challengeID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.ChallengeParticipantEntity{}).
		Where("challenge_id = ? AND user_id = ? AND status IN ?", challengeID, userID, []uint{uint(enum.StatusJoined), uint(enum.StatusPending)}).
		Count(&count).Error
	if err != nil {
		return false, exception.NewRepositoryError(err)
	}
	return count > 0, nil
}

// Helper functions
func toChallengeParticipantEntity(m *model.ChallengeParticipant) *entity.ChallengeParticipantEntity {
	return &entity.ChallengeParticipantEntity{
		ChallengeID: m.ChallengeID,
		UserID:      m.UserID,
		Status:      uint(m.Status),
	}
}

func toChallengeParticipantModel(e *entity.ChallengeParticipantEntity) *model.ChallengeParticipant {
	return &model.ChallengeParticipant{
		ID:          e.ID,
		ChallengeID: e.ChallengeID,
		UserID:      e.UserID,
		Status:      enum.ParticipantStatus(e.Status),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toChallengeParticipantModelWithID(e entity.ChallengeParticipantEntity, original *model.ChallengeParticipant) *model.ChallengeParticipant {
	original.ID = e.ID
	original.CreatedAt = e.CreatedAt
	original.UpdatedAt = e.UpdatedAt
	return original
}

func handleGetEntityWithCondition(err error, modelFunc func() *model.ChallengeParticipant) (*model.ChallengeParticipant, error) {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(err)
	}
	return modelFunc(), nil
}
