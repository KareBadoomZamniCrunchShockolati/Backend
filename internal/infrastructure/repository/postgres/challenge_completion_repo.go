package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"time"

	"gorm.io/gorm"
)

type ChallengeCompletionRepository struct {
	db *gorm.DB
}

func NewChallengeCompletionRepository(db *gorm.DB) *ChallengeCompletionRepository {
	return &ChallengeCompletionRepository{db: db}
}

func (r *ChallengeCompletionRepository) CreateCompletion(completion *model.ChallengeCompletion) (*model.ChallengeCompletion, error) {
	entity := &entity.ChallengeCompletionEntity{
		UserID:      completion.UserID,
		ChallengeID: completion.ChallengeID,
		IsCompleted: completion.IsCompleted,
		CompletedAt: completion.CompletedAt,
	}

	if err := r.db.Create(entity).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	completion.ID = entity.ID
	completion.CreatedAt = entity.CreatedAt
	completion.UpdatedAt = entity.UpdatedAt
	return completion, nil
}

func (r *ChallengeCompletionRepository) GetCompletion(userID, challengeID uint) (*model.ChallengeCompletion, error) {
	var e entity.ChallengeCompletionEntity
	err := r.db.Where("user_id = ? AND challenge_id = ?", userID, challengeID).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(err)
	}

	return toChallengeCompletionModel(&e), nil
}

func (r *ChallengeCompletionRepository) UpdateCompletion(completion *model.ChallengeCompletion) (*model.ChallengeCompletion, error) {
	entity := &entity.ChallengeCompletionEntity{
		UserID:      completion.UserID,
		ChallengeID: completion.ChallengeID,
		IsCompleted: completion.IsCompleted,
		CompletedAt: completion.CompletedAt,
	}
	entity.ID = completion.ID

	if err := r.db.Save(entity).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	completion.UpdatedAt = entity.UpdatedAt
	return completion, nil
}

func (r *ChallengeCompletionRepository) DeleteCompletion(userID, challengeID uint) error {
	result := r.db.Where("user_id = ? AND challenge_id = ?", userID, challengeID).
		Delete(&entity.ChallengeCompletionEntity{})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

func (r *ChallengeCompletionRepository) MarkChallengeCompleted(userID, challengeID uint) error {
	// Try to get existing completion record
	completion, err := r.GetCompletion(userID, challengeID)
	if err != nil {
		return err
	}

	if completion == nil {
		// Create new completion record
		completion = &model.ChallengeCompletion{
			UserID:      userID,
			ChallengeID: challengeID,
			IsCompleted: true,
			CompletedAt: time.Now(),
		}
		_, err = r.CreateCompletion(completion)
		return err
	}

	// Update existing record
	completion.IsCompleted = true
	completion.CompletedAt = time.Now()
	_, err = r.UpdateCompletion(completion)
	return err
}

func (r *ChallengeCompletionRepository) GetCompletedChallengesByUser(userID uint, offset, limit int) ([]*model.ChallengeCompletion, error) {
	var entities []entity.ChallengeCompletionEntity
	err := r.db.Where("user_id = ? AND is_completed = ?", userID, true).
		Offset(offset).Limit(limit).
		Order("completed_at DESC").
		Find(&entities).Error
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	completions := make([]*model.ChallengeCompletion, len(entities))
	for i, e := range entities {
		completions[i] = toChallengeCompletionModel(&e)
	}
	return completions, nil
}

func (r *ChallengeCompletionRepository) GetChallengesToProcess() ([]*model.ChallengeModel, error) {
	// Get challenges that have ended but not stopped
	now := time.Now()
	var challengeEntities []entity.ChallengeEntity

	err := r.db.Joins("JOIN challenge_participants cp ON cp.challenge_id = challenges.id").
		Where("challenges.end_time <= ? AND challenges.is_stopped = ? AND cp.status = ?",
			now, false, "joined").
		Group("challenges.id").
		Having("NOT EXISTS (SELECT 1 FROM challenge_completions cc WHERE cc.challenge_id = challenges.id AND cc.is_completed = ?)", true).
		Find(&challengeEntities).Error

	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	challenges := make([]*model.ChallengeModel, len(challengeEntities))
	for i, e := range challengeEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeCompletionRepository) GetCompletionRate(challengeID uint) (float64, error) {
	var totalParticipants, completedParticipants int64

	// Get total participants
	err := r.db.Model(&entity.ChallengeParticipantEntity{}).
		Where("challenge_id = ? AND status = ?", challengeID, "joined").
		Count(&totalParticipants).Error
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}

	if totalParticipants == 0 {
		return 0, nil
	}

	// Get completed participants
	err = r.db.Model(&entity.ChallengeCompletionEntity{}).
		Where("challenge_id = ? AND is_completed = ?", challengeID, true).
		Count(&completedParticipants).Error
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}

	return float64(completedParticipants) / float64(totalParticipants) * 100, nil
}

func (r *ChallengeCompletionRepository) GetUserCompletionStats(userID uint) (total, completed int, err error) {
	var totalChallenges, completedChallenges int64

	// Get total challenges user participated in
	err = r.db.Model(&entity.ChallengeParticipantEntity{}).
		Where("user_id = ? AND status = ?", userID, "joined").
		Count(&totalChallenges).Error
	if err != nil {
		return 0, 0, exception.NewRepositoryError(err)
	}

	// Get completed challenges
	err = r.db.Model(&entity.ChallengeCompletionEntity{}).
		Where("user_id = ? AND is_completed = ?", userID, true).
		Count(&completedChallenges).Error
	if err != nil {
		return 0, 0, exception.NewRepositoryError(err)
	}

	return int(totalChallenges), int(completedChallenges), nil
}

// Helper functions
func toChallengeCompletionModel(e *entity.ChallengeCompletionEntity) *model.ChallengeCompletion {
	return &model.ChallengeCompletion{
		ID:          e.ID,
		UserID:      e.UserID,
		ChallengeID: e.ChallengeID,
		IsCompleted: e.IsCompleted,
		CompletedAt: e.CompletedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
