package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/repository"
	"fmt"
	"time"
)

type ChallengeCompletionService struct {
	completionRepo  repository.ChallengeCompletionRepository
	challengeRepo   repository.ChallengeRepository
	participantRepo repository.ChallengeParticipantRepository
	userDayRepo     repository.UserDayRepository
	userRepo        repository.UserRepository
	categoryRepo    repository.CategoryRepository
}

func NewChallengeCompletionService(
	completionRepo repository.ChallengeCompletionRepository,
	challengeRepo repository.ChallengeRepository,
	participantRepo repository.ChallengeParticipantRepository,
	userDayRepo repository.UserDayRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
) *ChallengeCompletionService {
	return &ChallengeCompletionService{
		completionRepo:  completionRepo,
		challengeRepo:   challengeRepo,
		participantRepo: participantRepo,
		userDayRepo:     userDayRepo,
		userRepo:        userRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *ChallengeCompletionService) ProcessEndedChallenges() (int, error) {
	// challenges that have ended and need processing
	challenges, err := s.completionRepo.GetChallengesToProcess()
	// a list of challenges where 1.their end time has passed, 2. they have not been stopped
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}

	completedCount := 0
	for _, challenge := range challenges {
		// Get all participants for this challenge
		participants, err := s.participantRepo.GetParticipantsByChallenge(challenge.ID, 0, 1000)
		if err != nil {
			// Log error but continue with other challenges
			continue
		}

		for _, participant := range participants {
			if participant.Status != "joined" {
				continue
			}

			// Check if user achieved the goal in at least 80% of days
			isCompleted, _, err := s.checkUserGoalAchievement(
				participant.UserID,
				challenge.ID,
				challenge.StartTime,
				challenge.EndTime,
				challenge.Goal,
			)

			if err != nil {
				// Log error but continue with other participants
				continue
			}

			if isCompleted {
				// Mark challenge as completed for this user
				err = s.completionRepo.MarkChallengeCompleted(participant.UserID, challenge.ID)
				if err == nil {
					completedCount++
				}
			}
		}
	}

	return completedCount, nil
}

func (s *ChallengeCompletionService) checkUserGoalAchievement(
	userID uint,
	challengeID uint,
	startTime time.Time,
	endTime *time.Time,
	challengeGoal int,
) (bool, float64, error) {
	if endTime == nil {
		return false, 0, exception.NewBadRequestException(
			"CHALLENGE_NO_END_TIME",
			map[string]any{"challenge_id": challengeID},
		)
	}

	// Calculate total days in the challenge
	totalDays := int(endTime.Sub(startTime).Hours()/24) + 1

	if totalDays <= 0 {
		return false, 0, nil
	}

	// Get user's daily progress for this challenge
	progressData, _, err := s.userDayRepo.GetGoalProgressRange(
		userID,
		challengeID,
		startTime,
		*endTime,
	)
	if err != nil {
		return false, 0, exception.NewRepositoryError(err)
	}

	// Count days where goal was achieved
	successfulDays := 0
	totalProgress := 0

	for _, daily := range progressData {
		totalProgress += int(daily.Progress)
		if daily.Progress >= uint(challengeGoal) {
			successfulDays++
		}
	}

	// Calculate completion rate
	completionRate := float64(successfulDays) / float64(totalDays) * 100

	// User completed if they achieved goal in at least 80% of days
	isCompleted := completionRate >= 80.0

	return isCompleted, completionRate, nil
}

func (s *ChallengeCompletionService) CheckUserCompletion(userID, challengeID uint) (bool, float64, error) {
	// First check if we already have a completion record
	completion, err := s.completionRepo.GetCompletion(userID, challengeID)
	if err != nil {
		return false, 0, exception.NewRepositoryError(err)
	}

	if completion != nil && completion.IsCompleted {
		//  calculate the rate
		challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
		if err != nil {
			return true, 100, nil // Return true with 100% if we can't calculate
		}

		_, rate, err := s.checkUserGoalAchievement(
			userID,
			challengeID,
			challenge.StartTime,
			challenge.EndTime,
			challenge.Goal,
		)

		if err != nil {
			return true, 100, nil
		}

		return true, rate, nil
	}

	// Not completed yet
	return false, 0, nil
}

func (s *ChallengeCompletionService) GetUserCompletedChallenges(userID uint, offset, limit int) ([]*dto.CompletedChallengeDTO, error) {
	// Get completion records
	completions, err := s.completionRepo.GetCompletedChallengesByUser(userID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	result := make([]*dto.CompletedChallengeDTO, len(completions))

	for i, completion := range completions {
		// Get challenge details
		challenge, err := s.challengeRepo.GetChallengeByID(completion.ChallengeID, userID)
		if err != nil {
			// Skip for now
			continue
		}

		// Get category name
		categoryName := ""
		if challenge.CategoryID > 0 {
			category, err := s.categoryRepo.GetCategoryByID(challenge.CategoryID)
			if err == nil && category != nil {
				categoryName = category.Name
			}
		}

		// Get creator username
		creatorUsername := ""
		if challenge.CreatorID > 0 {
			user, err := s.userRepo.GetUserByID(challenge.CreatorID)
			if err == nil && user != nil {
				creatorUsername = user.Username
			}
		}

		// Calculate user's average progress
		avgProgress, err := s.calculateUserAverageProgress(userID, completion.ChallengeID)
		if err != nil {
			avgProgress = 0
		}

		// Get participant count
		participantCount, _ := s.participantRepo.GetParticipantCount(completion.ChallengeID)

		// Get like count
		likeCount, _ := s.getLikeCount(completion.ChallengeID)

		result[i] = &dto.CompletedChallengeDTO{
			ChallengePreviewDTO: dto.ChallengePreviewDTO{
				ID:                  challenge.ID,
				Title:               challenge.Title,
				Description:         challenge.Description,
				Rule:                challenge.Rule,
				CategoryName:        categoryName,
				CreatorUsername:     creatorUsername,
				CreatorID:           challenge.CreatorID,
				Visibility:          challenge.Visibility,
				Latitude:            challenge.Latitude,
				Longitude:           challenge.Longitude,
				Address:             challenge.Address,
				Goal:                challenge.Goal,
				MaxParticipants:     challenge.MaxParticipants,
				CurrentParticipants: participantCount,
				LikeCount:           likeCount,
				CommentCount:        0,
				StartTime:           challenge.StartTime,
				EndTime:             challenge.EndTime,
				Timezone:            challenge.Timezone,
				CreatedAt:           challenge.CreatedAt,
				IsUserParticipating: true,
				IsUserLiked:         false,
				CoverImage:          challenge.CoverImage,
			},
			CompletedAt:     completion.CompletedAt,
			ChallengeGoal:   challenge.Goal,
			AverageProgress: avgProgress,
		}

		// Calculate completion rate for this user
		_, completionRate, _ := s.checkUserGoalAchievement(
			userID,
			completion.ChallengeID,
			challenge.StartTime,
			challenge.EndTime,
			challenge.Goal,
		)
		result[i].CompletionRate = completionRate
	}

	return result, nil
}

func (s *ChallengeCompletionService) GetUserCompletionStats(userID uint) (*dto.UserCompletionStatsDTO, error) {
	total, completed, err := s.completionRepo.GetUserCompletionStats(userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	var completionRate float64
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	return &dto.UserCompletionStatsDTO{
		TotalChallenges:     total,
		CompletedChallenges: completed,
		CompletionRate:      completionRate,
	}, nil
}

func (s *ChallengeCompletionService) MarkChallengeAsCompleted(userID, challengeID uint) error {
	// Verify challenge exists and has ended
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}

	if challenge == nil {
		return exception.NewNotFoundException(
			"Challenge",
			fmt.Sprintf("%d", challengeID),
			"CHALLENGE_NOT_FOUND",
		)
	}

	// Check if challenge has ended
	if challenge.EndTime != nil && challenge.EndTime.After(time.Now()) {
		return exception.NewBadRequestException(
			"CHALLENGE_NOT_ENDED",
			map[string]any{
				"challenge_id": challengeID,
				"end_time":     challenge.EndTime,
			},
		)
	}

	// Mark as completed
	return s.completionRepo.MarkChallengeCompleted(userID, challengeID)
}

func (s *ChallengeCompletionService) GetChallengeCompletionRate(challengeID uint) (float64, error) {
	return s.completionRepo.GetCompletionRate(challengeID)
}

func (s *ChallengeCompletionService) calculateUserAverageProgress(userID, challengeID uint) (float64, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}

	if challenge.EndTime == nil {
		return 0, nil
	}

	progressData, _, err := s.userDayRepo.GetGoalProgressRange(
		userID,
		challengeID,
		challenge.StartTime,
		*challenge.EndTime,
	)
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}

	if len(progressData) == 0 {
		return 0, nil
	}

	total := 0
	for _, daily := range progressData {
		total += int(daily.Progress)
	}

	return float64(total) / float64(len(progressData)), nil
}

func (s *ChallengeCompletionService) getLikeCount(challengeID uint) (uint, error) {
	// For now, return 0, would be fully implemented after fixing likes bug
	return 0, nil
}
