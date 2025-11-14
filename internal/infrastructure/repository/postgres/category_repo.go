package postgres

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"fmt"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func ToCategoryModel(e *entity.ChallengeCategoryEntity) *model.ChallengeCategoryModel {
	if e == nil {
		return nil
	}
	return &model.ChallengeCategoryModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}

func ToCategoryResponseDTO(m *model.ChallengeCategoryModel) *dto.CategoryResponseDTO {
	if m == nil {
		return nil
	}
	return &dto.CategoryResponseDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	}
}

func ToCategoryResponseDTOs(models []*model.ChallengeCategoryModel) []*dto.CategoryResponseDTO {
	dtos := make([]*dto.CategoryResponseDTO, len(models))
	for i, m := range models {
		dtos[i] = ToCategoryResponseDTO(m)
	}
	return dtos
}

func (r *CategoryRepository) CreateCategory(category *model.ChallengeCategoryModel) (*model.ChallengeCategoryModel, error) {
	entity := &entity.ChallengeCategoryEntity{
		Name:        category.Name,
		Description: category.Description,
	}

	if err := r.DB.Create(entity).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	category.ID = entity.ID
	return category, nil
}

func (r *CategoryRepository) GetCategoryByID(id uint) (*model.ChallengeCategoryModel, error) {
	var entity entity.ChallengeCategoryEntity
	if err := r.DB.First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("Category", fmt.Sprintf("%d", id), "CATEGORY_NOT_FOUND")
		}
		return nil, exception.NewRepositoryError(err)
	}
	return ToCategoryModel(&entity), nil
}

func (r *CategoryRepository) GetCategoryByName(name string) (*model.ChallengeCategoryModel, error) {
	var entity entity.ChallengeCategoryEntity
	if err := r.DB.Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(err)
	}
	return ToCategoryModel(&entity), nil
}

func (r *CategoryRepository) GetAllCategories() ([]*model.ChallengeCategoryModel, error) {
	var entities []entity.ChallengeCategoryEntity
	if err := r.DB.Find(&entities).Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	result := make([]*model.ChallengeCategoryModel, len(entities))
	for i, e := range entities {
		result[i] = ToCategoryModel(&e)

	}
	return result, nil
}

func (r *CategoryRepository) UpdateCategory(category *model.ChallengeCategoryModel) (*model.ChallengeCategoryModel, error) {
	var entity entity.ChallengeCategoryEntity
	if err := r.DB.First(&entity, category.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundException("Category", fmt.Sprintf("%d", category.ID), "CATEGORY_NOT_FOUND")
		}
		return nil, exception.NewRepositoryError(err)
	}

	entity.Name = category.Name
	entity.Description = category.Description

	if err := r.DB.Save(&entity).Error; err != nil {
		return nil, exception.NewRepositoryUpdateError(err)
	}
	return category, nil
}

func (r *CategoryRepository) DeleteCategory(id uint) error {
	var entity entity.ChallengeCategoryEntity
	if err := r.DB.First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.NewNotFoundException("Category", string(id), "CATEGORY_NOT_FOUND")
		}
		return exception.NewRepositoryError(err)
	}
	if err := r.DB.Delete(&entity).Error; err != nil {
		return exception.NewRepositoryUpdateError(err)
	}
	return nil
}
