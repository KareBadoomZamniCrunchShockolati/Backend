package service

import (
	"strings"
	"time"

	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/repository"
)

type UserDayService struct {
	userDayRepo repository.UserDayRepository
}

func NewUserDayService(userDayRepo repository.UserDayRepository) *UserDayService {
	return &UserDayService{
		userDayRepo: userDayRepo,
	}
}

// Notes
func (s *UserDayService) CreateNote(userID, challengeID uint, date time.Time, note string) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}
	note = strings.TrimSpace(note)
	if note == "" {
		return exception.NewBadRequestException("NOTE_EMPTY", map[string]any{
			"field": "note",
		})
	}

	if err := s.userDayRepo.CreateNote(userID, challengeID, date, note); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) UpdateNote(userID, challengeID uint, date time.Time, note string) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}
	if err := s.userDayRepo.UpdateNote(userID, challengeID, date, note); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) GetNote(userID, challengeID uint, date time.Time) (string, error) {
	if userID == 0 {
		return "", exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return "", exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	note, err := s.userDayRepo.GetNote(userID, challengeID, date)
	if err != nil {
		return "", exception.NewRepositoryError(err)
	}
	return note, nil
}

func (s *UserDayService) DeleteNote(userID, challengeID uint, date time.Time) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	if err := s.userDayRepo.DeleteNote(userID, challengeID, date); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

// Feelings
func (s *UserDayService) CreateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}
	if feeling != enum.FeelingGood && feeling != enum.FeelingBad {
		return exception.NewBadRequestException("INVALID_FEELING", map[string]any{
			"field": "feeling",
			"got":   string(feeling),
			"allowed": []string{
				string(enum.FeelingGood),
				string(enum.FeelingBad),
			},
		})
	}

	if err := s.userDayRepo.CreateFeeling(userID, challengeID, date, feeling); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) UpdateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}
	if feeling != enum.FeelingGood && feeling != enum.FeelingBad {
		return exception.NewBadRequestException("INVALID_FEELING", map[string]any{
			"field": "feeling",
			"got":   string(feeling),
			"allowed": []string{
				string(enum.FeelingGood),
				string(enum.FeelingBad),
			},
		})
	}

	if err := s.userDayRepo.UpdateFeeling(userID, challengeID, date, feeling); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) GetFeeling(userID, challengeID uint, date time.Time) (enum.UserFeeling, error) {
	if userID == 0 {
		return "", exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return "", exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	feeling, err := s.userDayRepo.GetFeeling(userID, challengeID, date)
	if err != nil {
		return "", exception.NewRepositoryError(err)
	}
	return feeling, nil
}

func (s *UserDayService) DeleteFeeling(userID, challengeID uint, date time.Time) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	if err := s.userDayRepo.DeleteFeeling(userID, challengeID, date); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

// Goal
func (s *UserDayService) CreateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	if err := s.userDayRepo.CreateGoalProgress(userID, challengeID, date, progress); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) UpdateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	if progress == 0 {
		return exception.NewBadRequestException("INVALID_PROGRESS", map[string]any{
			"field": "progress",
			"min":   1,
			"got":   progress,
		})
	}

	if err := s.userDayRepo.UpdateGoalProgress(userID, challengeID, date, progress); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) GetGoalProgress(userID, challengeID uint, date time.Time) (uint, error) {
	if userID == 0 {
		return 0, exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return 0, exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	progress, err := s.userDayRepo.GetGoalProgress(userID, challengeID, date)
	if err != nil {
		return 0, exception.NewRepositoryError(err)
	}
	return progress, nil
}

func (s *UserDayService) DeleteGoalProgress(userID, challengeID uint, date time.Time) error {
	if userID == 0 {
		return exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	if err := s.userDayRepo.DeleteGoalProgress(userID, challengeID, date); err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserDayService) GetGoalProgressChart(userID, challengeID uint, start, end time.Time) (*dto.DaysProgressResponse, error) {
	if userID == 0 {
		return nil, exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return nil, exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}
	if start.After(end) {
		return nil, exception.NewBadRequestException("INVALID_DATE_RANGE", map[string]any{
			"start": start.Format("2006-01-02"),
			"end":   end.Format("2006-01-02"),
		})
	}

	progress, mean, err := s.userDayRepo.GetGoalProgressRange(userID, challengeID, start, end)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
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
		return nil, exception.NewBadRequestException("USER_ID_REQUIRED", map[string]any{
			"field": "userID",
		})
	}
	if challengeID == 0 {
		return nil, exception.NewBadRequestException("CHALLENGE_ID_REQUIRED", map[string]any{
			"field": "challengeID",
		})
	}

	good, bad, err := s.userDayRepo.GetFeelingCounts(userID, challengeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	return &dto.FeelingCountResponse{
		GoodDays: good,
		BadDays:  bad,
	}, nil
}