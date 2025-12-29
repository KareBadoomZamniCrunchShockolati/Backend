package repository

import "challenge-app/internal/domain/model"

type MedalRepository interface {
	// List medals for a category
	GetMedalsByCategory(categoryID uint) ([]*model.MedalModel, error)
}
