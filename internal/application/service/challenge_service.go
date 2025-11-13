package service

import (
	"fmt"
	"time"

	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
)

type ChallengeService struct {
	ChallengeRepo repository.ChallengeRepository
	UserRepo      repository.UserRepository
	CategoryRepo  repository.CategoryRepository
}

func NewChallengeService(
	chRepo repository.ChallengeRepository, userRepo repository.UserRepository, categoryRepo repository.CategoryRepository,
) *ChallengeService {
	return &ChallengeService{
		ChallengeRepo: chRepo,
		UserRepo:      userRepo,
		CategoryRepo:  categoryRepo,
	}
}


func toChallengeModelFromCreateDTO(input *dto.CreateChallengeDTO) *model.ChallengeModel {
	if input == nil {
		return nil
	}

	return &model.ChallengeModel{
		Title:           input.Title,
		Description:     input.Description,
		CategoryID:      input.CategoryID,
		CreatorID:       input.CreatorID,
		MaxParticipants: input.MaxParticipants,
		Visibility:      input.Visibility,
		Rule:            input.Rule,
		StartTime:       input.StartTime,
		EndTime:         &input.EndTime,
		Timezone:        input.Timezone,
		ImageURL:        input.ImageURL,
		IsStopped:       false,
	}
}

func applyUpdateDTOToChallenge(existing *model.ChallengeModel, input *dto.UpdateChallengeDTO) {
	if input == nil || existing == nil {
		return
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.CategoryID != nil {
		existing.CategoryID = *input.CategoryID
	}
	if input.MaxParticipants != nil {
		existing.MaxParticipants = *input.MaxParticipants
	}
	if input.Visibility != nil {
		existing.Visibility = *input.Visibility
	}
	if input.Rule != nil {
		existing.Rule = *input.Rule
	}
	if input.ImageURL != nil {
		existing.ImageURL = *input.ImageURL
	}
	if input.Timezone != nil {
		existing.Timezone = *input.Timezone
	}
	if input.StartTime != nil && existing.StartTime.After(*input.StartTime) {
		existing.StartTime = *input.StartTime
	}
	if input.EndTime != nil {
		existing.EndTime = input.EndTime
	}
	if input.CommentsEnabled != nil {
		existing.CommentsEnabled = *input.CommentsEnabled
	}
	if input.IsStopped != nil {
		existing.IsStopped = *input.IsStopped
	}
}


func (s *ChallengeService) ToChallengeResponseDTO(ch *model.ChallengeModel, currentUserID uint) (*dto.ChallengeResponseDTO, error) {
	creator, err := s.UserRepo.GetUserByID(ch.CreatorID)
	if err != nil {
		return nil, err
	}

	category, err := s.CategoryRepo.GetCategoryByID(ch.CategoryID)
	if err != nil {
		return nil, err
	}

	// count, err := s.ChallengeRepo.GetParticipantCount(ch.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return &dto.ChallengeResponseDTO{
		ID:                  ch.ID,
		Title:               ch.Title,
		Description:         ch.Description,
		CategoryName:        category.Name,
		CreatorUsername:     creator.Username,
		MaxParticipants:     ch.MaxParticipants,
		CurrentParticipants: 0,//implement get participant count later and fix this 
		Visibility:          ch.Visibility,
		ImageURL:            ch.ImageURL,
		Rule:                ch.Rule,
		Timezone:            ch.Timezone,
		StartTime:           ch.StartTime,
		EndTime:             ch.EndTime,
		IsStopped:           ch.IsStopped,
		CommentsEnabled:     ch.CommentsEnabled,
		CreatedAt:           ch.CreatedAt,
	}, nil
}

func (s *ChallengeService) ToChallengeResponseDTOs(challenges []*model.ChallengeModel, currentUserID uint) ([]*dto.ChallengeResponseDTO, error) {
	result := make([]*dto.ChallengeResponseDTO, 0, len(challenges))
	for _, ch := range challenges {
		dto, err := s.ToChallengeResponseDTO(ch, currentUserID)
		if err != nil {
			return nil, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (s *ChallengeService) CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("Input cannot be nil", "CHALLENGE_CREATE_BAD_INPUT", nil)
	}

	ch := toChallengeModelFromCreateDTO(input)

	created, err := s.ChallengeRepo.CreateChallenge(ch)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if created == nil {
		return nil, exception.NewInternalServerException("ChallengeRepo.CreateChallenge returned nil", "CHALLENGE_CREATE_FAILED", nil)
	}

	return created, nil
}

func (s *ChallengeService) GetChallengeByID(id uint) (*model.ChallengeModel, error) {
	ch, err := s.ChallengeRepo.GetChallengeByID(id)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			return nil, err
		}
		return nil, exception.NewRepositoryError(err)
	}
	if ch == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}

	return ch, nil
}

func (s *ChallengeService) GetAllChallenges() ([]*model.ChallengeModel, error) {
	challenges, err := s.ChallengeRepo.GetAllChallenges()
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return challenges, nil
}

// UpdateChallenge updates a challenge
func (s *ChallengeService) UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error) {
	existing, err := s.ChallengeRepo.GetChallengeByID(id)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if existing == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}

	if existing.CreatorID != currentUserID {
		return nil, exception.NewUnauthorizedException("Only the creator can update this challenge", "CHALLENGE_UPDATE_FORBIDDEN")
	}

	applyUpdateDTOToChallenge(existing, input)

	if input.StartTime != nil && existing.StartTime.Before(time.Now()) {
		return nil, exception.NewBadRequestException("Cannot change StartTime, challenge already started", "CHALLENGE_ALREADY_STARTED", nil)
	}

	updated, err := s.ChallengeRepo.UpdateChallenge(existing)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if updated == nil {
		return nil, exception.NewInternalServerException(fmt.Sprintf("ChallengeRepo.UpdateChallenge returned nil for ID %d", id), "CHALLENGE_UPDATE_FAILED", nil)
	}

	return updated, nil
}

// DeleteChallenge deletes a challenge
func (s *ChallengeService) DeleteChallenge(id uint, currentUserID uint) error {
	ch, err := s.ChallengeRepo.GetChallengeByID(id)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if ch == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}

	if ch.CreatorID != currentUserID {
		return exception.NewUnauthorizedException("Only the creator can delete this challenge", "CHALLENGE_DELETE_FORBIDDEN")
	}

	if err := s.ChallengeRepo.DeleteChallenge(id); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}
