package postgres

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"time"

	"gorm.io/gorm"
)

type UserDayRepository struct {
	db *gorm.DB
}

func NewUserDayRepository(db *gorm.DB) *UserDayRepository {
	return &UserDayRepository{db: db}
}

func (r *UserDayRepository) normalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// Notes 
func (r *UserDayRepository) CreateNote(userID, challengeID uint, date time.Time, note string) error {
	ent := entity.UserDailyNoteEntity{
		UserID:      userID,
		ChallengeID: challengeID,
		Date:        r.normalizeDate(date),
		Note:        note,
	}
	if err := r.db.Create(&ent).Error; err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (r *UserDayRepository) UpdateNote(userID, challengeID uint, date time.Time, note string) error {
	result := r.db.Model(&entity.UserDailyNoteEntity{}).
		Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		Updates(map[string]interface{}{"note": note})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundException("DailyNote", "note", "USER_DAY_NOTE_NOT_FOUND")
	}
	return nil
}

func (r *UserDayRepository) GetNote(userID, challengeID uint, date time.Time) (string, error) {
	var ent entity.UserDailyNoteEntity
	err := r.db.Select("note").
		Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		First(&ent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil 
		}
		return "", exception.NewRepositoryError(err)
	}
	return ent.Note, nil
}

func (r *UserDayRepository) DeleteNote(userID, challengeID uint, date time.Time) error {
	result := r.db.Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		Delete(&entity.UserDailyNoteEntity{})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

// Feelings 
func (r *UserDayRepository) CreateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error {
	ent := entity.UserDailyFeelingEntity{
		UserID:      userID,
		ChallengeID: challengeID,
		Date:        r.normalizeDate(date),
		Feeling:     feeling,
	}
	if err := r.db.Create(&ent).Error; err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (r *UserDayRepository) UpdateFeeling(userID, challengeID uint, date time.Time, feeling enum.UserFeeling) error {
	result := r.db.Model(&entity.UserDailyFeelingEntity{}).
		Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		Updates(map[string]interface{}{"feeling": feeling})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundException("DailyFeeling", "feeling", "USER_DAY_FEELING_NOT_FOUND")
	}
	return nil
}

func (r *UserDayRepository) GetFeeling(userID, challengeID uint, date time.Time) (enum.UserFeeling, error) {
	var ent entity.UserDailyFeelingEntity
	err := r.db.Select("feeling").
		Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		First(&ent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", exception.NewRepositoryError(err)
	}
	return ent.Feeling, nil
}

func (r *UserDayRepository) DeleteFeeling(userID, challengeID uint, date time.Time) error {
	result := r.db.Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		Delete(&entity.UserDailyFeelingEntity{})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

// Goal Progress 
func (r *UserDayRepository) CreateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error {
	ent := entity.UserDailyGoalProgressEntity{
		UserID:      userID,
		ChallengeID: challengeID,
		Date:        r.normalizeDate(date),
		Progress:    progress,
	}
	if err := r.db.Create(&ent).Error; err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (r *UserDayRepository) UpdateGoalProgress(userID, challengeID uint, date time.Time, progress uint) error {
	result := r.db.Model(&entity.UserDailyGoalProgressEntity{}).
		Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		Updates(map[string]interface{}{"progress": progress})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundException("DailyGoalProgress", "progress", "USER_DAY_PROGRESS_NOT_FOUND")
	}
	return nil
}

func (r *UserDayRepository) GetGoalProgress(userID, challengeID uint, date time.Time) (uint, error) {
	var ent entity.UserDailyGoalProgressEntity
	err := r.db.Select("progress").
		Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		First(&ent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, exception.NewRepositoryError(err)
	}
	return ent.Progress, nil
}

func (r *UserDayRepository) DeleteGoalProgress(userID, challengeID uint, date time.Time) error {
	result := r.db.Where("user_id = ? AND challenge_id = ? AND date = ?", userID, challengeID, r.normalizeDate(date)).
		Delete(&entity.UserDailyGoalProgressEntity{})
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

func (r *UserDayRepository) GetFeelingCounts(userID, challengeID uint) (good, bad int, err error) {
	type Result struct {
		Feeling string
		Count   int
	}
	var results []Result
	queryErr := r.db.Table("user_daily_feelings").
		Select("feeling, COUNT(*) as count").
		Where("user_id = ? AND challenge_id = ?", userID, challengeID).
		Group("feeling").
		Find(&results).Error
	if queryErr != nil {
		return 0, 0, exception.NewRepositoryError(queryErr)
	}

	for _, res := range results {
		if enum.UserFeeling(res.Feeling) == enum.FeelingGood {
			good = res.Count
		} else if enum.UserFeeling(res.Feeling) == enum.FeelingBad {
			bad = res.Count
		}
	}
	return good, bad, nil
}

func (r *UserDayRepository) GetGoalProgressRange(userID, challengeID uint, start, end time.Time) ([]dto.DailyProgress, float64, error) {
	start = r.normalizeDate(start)
	end = r.normalizeDate(end)

	type Row struct {
		Date     time.Time `gorm:"column:date"`
		Progress uint      `gorm:"column:progress"`
	}
	var rows []Row
	err := r.db.Table("user_daily_goal_progress").
		Select("date, progress").
		Where("user_id = ? AND challenge_id = ? AND date BETWEEN ? AND ?", userID, challengeID, start, end).
		Order("date").
		Find(&rows).Error
	if err != nil {
		return nil, 0, exception.NewRepositoryError(err)
	}

	allDates := []string{}
	progressMap := make(map[string]uint)
	current := start
	for !current.After(end) {
		d := current.Format("2006-01-02")
		allDates = append(allDates, d)
		progressMap[d] = 0
		current = current.Add(24 * time.Hour)
	}

	for _, row := range rows {
		d := row.Date.Format("2006-01-02")
		progressMap[d] = row.Progress
	}

	var progressList []uint
	sum := uint(0)
	for _, d := range allDates {
		p := progressMap[d]
		progressList = append(progressList, p)
		sum += p
	}

	mean := 0.0
	if len(allDates) > 0 {
		mean = float64(sum) / float64(len(allDates))
	}

	var result []dto.DailyProgress
	for i, d := range allDates {
		result = append(result, dto.DailyProgress{Date: d, Progress: progressList[i]})
	}
	return result, mean, nil
}
