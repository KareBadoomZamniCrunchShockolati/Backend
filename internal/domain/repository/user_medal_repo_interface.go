package repository

import "challenge-app/internal/domain/model"

type UserMedalRepository interface {
	// award a medal record to a user
	AwardMedal(userID, categoryID uint, medalType string) (*model.UserMedalModel, error)

	GetUserMedals(userID uint) ([]*model.UserMedalModel, error)

	// increment and return user's completion count for a category
	IncrementUserCategoryCount(userID, categoryID uint) (int, error)

	GetUserCategoryCount(userID, categoryID uint) (int, error)

	//  get or set selected medals. up to 3 can be selected.
	GetSelectedMedals(userID uint) ([]string, error)
	SetSelectedMedals(userID uint, medals []string) error
}
