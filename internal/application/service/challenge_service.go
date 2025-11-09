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
}

func NewChallengeService(chRepo repository.ChallengeRepository) *ChallengeService {
	return &ChallengeService{
		ChallengeRepo: chRepo,
	}
}

func toChallengeModelFromCreateDTO(input *dto.CreateChallengeDTO) *model.ChallengeModel {
	if input == nil {
		return nil
	}

	return &model.ChallengeModel{
		Title:           input.Title,
		Description:     input.Description,
		Category:        input.Category,
		CreatorID:       input.CreatorID,
		MaxParticipants: input.MaxParticipants,
		Visibility:      input.Visibility,
		Rule:            input.Rule,
		StartTime:       input.StartTime,
		EndTime:         &input.EndTime,
		Timezone:        input.Timezone,
		ImageURL:        input.ImageURL,
		Stopped:         false,
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
	if input.Category != nil {
		existing.Category = *input.Category
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
}


func (s *ChallengeService) CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("Input cannot be nil", "CHALLENGE_CREATE_BAD_INPUT", nil)
	}

	ch := toChallengeModelFromCreateDTO(input)

	created, err := s.ChallengeRepo.CreateChallenge(ch)
	if err != nil {
		panic(exception.NewRepositoryError("Failed to create challenge", err))
	}
	if created == nil {
		panic(exception.NewInternalServerException("ChallengeRepo.CreateChallenge returned nil", "CHALLENGE_CREATE_FAILED", nil))
	}

	return created, nil
}

func (s *ChallengeService) GetChallengeByID(id uint) (*model.ChallengeModel, error) {
	ch, err := s.ChallengeRepo.GetChallengeByID(id)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			return nil, err
		}
		panic(exception.NewRepositoryError(fmt.Sprintf("GetChallengeByID %d failed", id), err))
	}
	if ch == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}

	return ch, nil
}

// UpdateChallenge updates a challenge with partial update DTO
func (s *ChallengeService) UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error) {
	existing, err := s.ChallengeRepo.GetChallengeByID(id)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			return nil, err
		}
		panic(exception.NewRepositoryError(fmt.Sprintf("GetChallengeByID %d failed", id), err))
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
		panic(exception.NewRepositoryError(fmt.Sprintf("UpdateChallenge %d failed", id), err))
	}
	if updated == nil {
		panic(exception.NewInternalServerException(fmt.Sprintf("ChallengeRepo.UpdateChallenge returned nil for ID %d", id), "CHALLENGE_UPDATE_FAILED", nil))
	}

	return updated, nil
}

// DeleteChallenge deletes a challenge
func (s *ChallengeService) DeleteChallenge(id uint, currentUserID uint) error {
	ch, err := s.ChallengeRepo.GetChallengeByID(id)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			return err
		}
		panic(exception.NewRepositoryError(fmt.Sprintf("GetChallengeByID %d failed", id), err))
	}

	if ch.CreatorID != currentUserID {
		return exception.NewUnauthorizedException("Only the creator can delete this challenge", "CHALLENGE_DELETE_FORBIDDEN")
	}

	if err := s.ChallengeRepo.DeleteChallenge(id); err != nil {
		panic(exception.NewRepositoryError(fmt.Sprintf("DeleteChallenge %d failed", id), err))
	}
	return nil
}