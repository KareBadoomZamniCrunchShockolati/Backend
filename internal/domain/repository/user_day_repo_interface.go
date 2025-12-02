package repository

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"time"
)

type UserDayRepository interface {
	CreateNote(userID, challengeID uint, date time.Time, note string) error
	UpdateNote(userID, challengeID uint, date time.Time, note string) error
	GetNote(userID, challengeID uint, date time.Time) (string, error)
	DeleteNote(userID, challengeID uint, date time.Time) error

	CreateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error
	UpdateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error
	GetFeeling(userID, challengeID uint, date time.Time) (enum.UserFeeling, error)
	DeleteFeeling(userID, challengeID uint, date time.Time) error

	CreateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error
	UpdateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error
	GetGoalProgress(userID, challengeID uint, date time.Time) (uint, error)
	DeleteGoalProgress(userID, challengeID uint, date time.Time) error

	GetGoalProgressRange(userID, challengeID uint, start, end time.Time) ([]dto.DailyProgress, float64, error)
	GetFeelingCounts(userID, challengeID uint) (good, bad int, err error)
}