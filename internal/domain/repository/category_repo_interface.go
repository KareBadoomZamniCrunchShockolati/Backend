package repository

import "challenge-app/internal/domain/model"

type CategoryRepository interface {
	CreateCategory(category *model.ChallengeCategoryModel) (*model.ChallengeCategoryModel, error)
	GetCategoryByID(id uint) (*model.ChallengeCategoryModel, error)
	GetCategoryByName(name string) (*model.ChallengeCategoryModel, error)
	GetAllCategories() ([]*model.ChallengeCategoryModel, error)
	UpdateCategory(category *model.ChallengeCategoryModel) (*model.ChallengeCategoryModel, error)
	DeleteCategory(id uint) error
}
