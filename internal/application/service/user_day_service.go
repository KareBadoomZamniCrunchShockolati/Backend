package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/repository"
	"time"
)

type UserDayService struct {
	userDayRepo repository.UserDayRepository
}

func NewUserDayService(userDayRepo repository.UserDayRepository,) *UserDayService {
	return &UserDayService{
		userDayRepo: userDayRepo,
	}
}

// Notes
func (s *UserDayService) CreateNote(userID, challengeID uint, date time.Time, note string) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}
	if note == "" {
		return exception.NewBadRequestException("Note cannot be empty", "NOTE_EMPTY", nil)
	}

	return s.userDayRepo.CreateNote(userID, challengeID, date, note)
}

func (s *UserDayService) UpdateNote(userID, challengeID uint, date time.Time, note string) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	err := s.userDayRepo.UpdateNote(userID, challengeID, date, note)
	if err != nil {
		return err 
	}
	return nil
}

func (s *UserDayService) GetNote(userID, challengeID uint, date time.Time) (string, error) {
	if userID == 0 {
		return "", exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return "", exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.GetNote(userID, challengeID, date)
}

func (s *UserDayService) DeleteNote(userID, challengeID uint, date time.Time) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.DeleteNote(userID, challengeID, date)
}

// Feelings
func (s *UserDayService) CreateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}
	if feeling != enum.FeelingGood && feeling != enum.FeelingBad {
		return exception.NewBadRequestException("Invalid feeling value", "INVALID_FEELING", nil)
	}

	return s.userDayRepo.CreateFeeling(userID, challengeID, date, feeling)
}

func (s *UserDayService) UpdateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}
	if feeling != enum.FeelingGood && feeling != enum.FeelingBad {
		return exception.NewBadRequestException("Invalid feeling value", "INVALID_FEELING", nil)
	}

	err := s.userDayRepo.UpdateFeeling(userID, challengeID, date, feeling)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserDayService) GetFeeling(userID, challengeID uint, date time.Time) (enum.UserFeeling, error) {
	if userID == 0 {
		return "", exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return "", exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.GetFeeling(userID, challengeID, date)
}

func (s *UserDayService) DeleteFeeling(userID, challengeID uint, date time.Time) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.DeleteFeeling(userID, challengeID, date)
}

// Goal
func (s *UserDayService) CreateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.CreateGoalProgress(userID, challengeID, date, progress)
}

func (s *UserDayService) UpdateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	err := s.userDayRepo.UpdateGoalProgress(userID, challengeID, date, progress)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserDayService) GetGoalProgress(userID, challengeID uint, date time.Time) (uint, error) {
	if userID == 0 {
		return 0, exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return 0, exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.GetGoalProgress(userID, challengeID, date)
}

func (s *UserDayService) DeleteGoalProgress(userID, challengeID uint, date time.Time) error {
	if userID == 0 {
		return exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	return s.userDayRepo.DeleteGoalProgress(userID, challengeID, date)
}

func (s *UserDayService) GetGoalProgressChart(userID, challengeID uint, start, end time.Time) (*dto.DaysProgressResponse, error) {
	if userID == 0 {
		return nil, exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return nil, exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}
	if start.After(end) {
		return nil, exception.NewBadRequestException("Start date must be before end date", "INVALID_DATE_RANGE", nil)
	}

	progress, mean, err := s.userDayRepo.GetGoalProgressRange(userID, challengeID, start, end)
	if err != nil {
		return nil, err
	}

	dates := make([]string, len(progress))
	values := make([]uint, len(progress))
	for i, p := range progress {
		dates[i] = p.Date
		values[i] = p.Progress
	}

	return &dto.DaysProgressResponse{
		Dates:        dates,
		Progress:     values,
		MeanProgress: mean,
	}, nil
}

func (s *UserDayService) GetTotalFeelings(userID, challengeID uint) (*dto.FeelingCountResponse, error) {
	if userID == 0 {
		return nil, exception.NewBadRequestException("User ID is required", "USER_ID_REQUIRED", nil)
	}
	if challengeID == 0 {
		return nil, exception.NewBadRequestException("Challenge ID is required", "CHALLENGE_ID_REQUIRED", nil)
	}

	good, bad, err := s.userDayRepo.GetFeelingCounts(userID, challengeID)
	if err != nil {
		return nil, err
	}

	return &dto.FeelingCountResponse{
		GoodDays: good,
		BadDays:  bad,
	}, nil
}