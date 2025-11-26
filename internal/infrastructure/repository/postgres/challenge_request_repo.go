// internal/infrastructure/repository/postgres/challenge_join_request_repo.go
package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type ChallengeRequestRepository struct {
	db *gorm.DB
}

func NewChallengeRequestRepository(db *gorm.DB) *ChallengeRequestRepository {
	return &ChallengeRequestRepository{db: db}
}

func (r *ChallengeRequestRepository) CreateJoinRequest(request *model.ChallengeRequest) (*model.ChallengeRequest, error) {
	entity := &entity.ChallengeRequestEntity{
		ChallengeID: request.ChallengeID,
		UserID:      request.RequesterID,
		Status:      uint(request.Status),
	}
	result := r.db.Create(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	request.ID = entity.ID
	request.CreatedAt = entity.CreatedAt
	request.UpdatedAt = entity.UpdatedAt
	return request, nil
}

func (r *ChallengeRequestRepository) GetJoinRequest(requestID uint) (*model.ChallengeRequest, error) {
	var entity entity.ChallengeRequestEntity
	result := r.db.First(&entity, requestID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(result.Error)
	}
	return &model.ChallengeRequest{
		ID:          entity.ID,
		ChallengeID: entity.ChallengeID,
		RequesterID: entity.UserID,
		Status:      enum.RequestStatus(entity.Status),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func (r *ChallengeRequestRepository) GetChallengeRequests(challengeID uint) ([]*model.ChallengeRequest, error) {
	var entities []entity.ChallengeRequestEntity
	result := r.db.Where("challenge_id = ?", challengeID).Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	requests := make([]*model.ChallengeRequest, len(entities))
	for i, e := range entities {
		requests[i] = &model.ChallengeRequest{
			ID:          e.ID,
			ChallengeID: e.ChallengeID,
			RequesterID: e.UserID,
			Status:      enum.RequestStatus(e.Status),
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		}
	}
	return requests, nil
}

func (r *ChallengeRequestRepository) GetUserJoinRequests(userID uint) ([]*model.ChallengeRequest, error) {
	var entities []entity.ChallengeRequestEntity
	result := r.db.Where("user_id = ?", userID).Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	requests := make([]*model.ChallengeRequest, len(entities))
	for i, e := range entities {
		requests[i] = &model.ChallengeRequest{
			ID:          e.ID,
			ChallengeID: e.ChallengeID,
			RequesterID: e.UserID,
			Status:      enum.RequestStatus(e.Status),
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		}
	}
	return requests, nil
}

func (r *ChallengeRequestRepository) UpdateJoinRequest(request *model.ChallengeRequest) (*model.ChallengeRequest, error) {
	entity := &entity.ChallengeRequestEntity{
		ChallengeID: request.ChallengeID,
		UserID:      request.RequesterID,
		Status:      uint(request.Status),
	}

	result := r.db.Save(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	request.UpdatedAt = entity.UpdatedAt
	return request, nil
}

func (r *ChallengeRequestRepository) DeleteJoinRequest(requestID uint) error {
	result := r.db.Delete(&entity.ChallengeRequestEntity{}, requestID)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}
