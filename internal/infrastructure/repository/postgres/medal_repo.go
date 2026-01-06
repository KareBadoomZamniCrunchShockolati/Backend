package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MedalRepository struct {
	db *gorm.DB
}

func NewMedalRepository(db *gorm.DB) *MedalRepository {
	return &MedalRepository{db: db}
}

func (r *MedalRepository) GetMedalsByCategory(categoryID uint) ([]*model.MedalModel, error) {
	var entities []entity.MedalEntity
	if err := r.db.Where("category_id = ?", categoryID).Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	res := make([]*model.MedalModel, len(entities))
	for i, e := range entities {
		res[i] = &model.MedalModel{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			CategoryID:  e.CategoryID,
			Type:        e.Type,
		}
	}
	return res, nil
}

// User medal related repo
type UserMedalRepositoryImpl struct {
	db *gorm.DB
}

func NewUserMedalRepository(db *gorm.DB) *UserMedalRepositoryImpl {
	return &UserMedalRepositoryImpl{db: db}
}

func (r *UserMedalRepositoryImpl) AwardMedal(userID, categoryID uint, medalType string) (*model.UserMedalModel, error) {
	e := &entity.UserMedalEntity{
		UserID:     userID,
		CategoryID: categoryID,
		Type:       medalType,
		AwardedAt:  time.Now(),
	}
	res := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(e)
	if res.Error != nil {
		return nil, exception.NewRepositoryError(res.Error)
	}

	if res.RowsAffected == 0 {
		// Medal already exists, return existing
		var exist entity.UserMedalEntity
		if err := r.db.Where("user_id = ? AND category_id = ? AND type = ?", userID, categoryID, medalType).First(&exist).Error; err != nil {
			return nil, exception.NewRepositoryError(err)
		}
		return &model.UserMedalModel{ID: exist.ID, UserID: exist.UserID, CategoryID: exist.CategoryID, Type: exist.Type, AwardedAt: exist.AwardedAt}, nil
	}

	return &model.UserMedalModel{ID: e.ID, UserID: e.UserID, CategoryID: e.CategoryID, Type: e.Type, AwardedAt: e.AwardedAt}, nil
}

func (r *UserMedalRepositoryImpl) GetUserMedals(userID uint) ([]*model.UserMedalModel, error) {
	var entities []entity.UserMedalEntity
	if err := r.db.Where("user_id = ?", userID).Order("awarded_at DESC").Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	res := make([]*model.UserMedalModel, len(entities))
	for i, e := range entities {
		res[i] = &model.UserMedalModel{ID: e.ID, UserID: e.UserID, CategoryID: e.CategoryID, Type: e.Type, AwardedAt: e.AwardedAt}
	}
	return res, nil
}

func (r *UserMedalRepositoryImpl) IncrementUserCategoryCount(userID, categoryID uint) (int, error) {
	// Use atomic upsert to increment count and return the new count
	var newCount int
	raw := `INSERT INTO user_category_progress_entities (user_id, category_id, count, updated_at)
	VALUES (?, ?, 1, now())
	ON CONFLICT (user_id, category_id) DO UPDATE SET count = user_category_progress_entities.count + 1, updated_at = now()
	RETURNING count;`

	if err := r.db.Raw(raw, userID, categoryID).Scan(&newCount).Error; err != nil {
		return 0, exception.NewRepositoryError(err)
	}
	return newCount, nil
}

func (r *UserMedalRepositoryImpl) GetUserCategoryCount(userID, categoryID uint) (int, error) {
	var p entity.UserCategoryProgressEntity
	if err := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, exception.NewRepositoryError(err)
	}
	return p.Count, nil
}

func (r *UserMedalRepositoryImpl) GetSelectedMedals(userID uint) ([]string, error) {
	var e entity.UserSelectedMedalsEntity
	if err := r.db.Where("user_id = ?", userID).First(&e).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		return nil, exception.NewRepositoryError(err)
	}
	var medals []string
	if err := json.Unmarshal([]byte(e.Medals), &medals); err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return medals, nil
}

func (r *UserMedalRepositoryImpl) SetSelectedMedals(userID uint, medals []string) error {
	// enforce max 3 at repository level
	if len(medals) > 3 {
		medals = medals[:3]
	}
	b, _ := json.Marshal(medals)
	var e entity.UserSelectedMedalsEntity
	tx := r.db.Where("user_id = ?", userID).First(&e)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			ne := entity.UserSelectedMedalsEntity{UserID: userID, Medals: string(b), UpdatedAt: time.Now()}
			if err := r.db.Create(&ne).Error; err != nil {
				return exception.NewRepositoryError(err)
			}
			return nil
		}
		return exception.NewRepositoryError(tx.Error)
	}
	e.Medals = string(b)
	e.UpdatedAt = time.Now()
	if err := r.db.Save(&e).Error; err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}
