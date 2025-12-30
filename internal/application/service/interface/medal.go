package serviceinterface

import "challenge-app/internal/application/dto"

type MedalServicer interface {
	OnCategoryCompleted(userID, categoryID uint) error
	GetUserMedals(userID uint) ([]*dto.MedalDTO, error)
	GetSelectedMedals(userID uint) ([]string, error)
	SelectMedals(userID uint, medals []string) error
	DeleteSelectedMedals(userID uint, medals []string) error
}
