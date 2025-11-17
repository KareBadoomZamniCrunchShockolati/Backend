package serviceinterface

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/model"
)

type CategoryServicer interface {
	CreateCategory(dto *dto.CreateCategoryDTO) (*model.ChallengeCategoryModel, error)
	GetCategoryByID(id uint) (*model.ChallengeCategoryModel, error)
	GetAllCategories() ([]*model.ChallengeCategoryModel, error)
	UpdateCategory(id uint, dto *dto.UpdateCategoryDTO) (*model.ChallengeCategoryModel, error)
	DeleteCategory(id uint) error
}
