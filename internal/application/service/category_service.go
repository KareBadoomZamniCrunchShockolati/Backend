package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: repo}
}

func (s *CategoryService) CreateCategory(dto *dto.CreateCategoryDTO) (*model.ChallengeCategoryModel, error) {
	existing, err := s.categoryRepo.GetCategoryByName(dto.Name)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if existing != nil {
		return nil, exception.NewConflictException("CATEGORY_ALREADY_EXISTS", "Category", "name", dto.Name)
	}

	return s.categoryRepo.CreateCategory(&model.ChallengeCategoryModel{
		Name:        dto.Name,
		Description: dto.Description,
	})
}

func (s *CategoryService) GetCategoryByID(id uint) (*model.ChallengeCategoryModel, error) {
	return s.categoryRepo.GetCategoryByID(id)
}

func (s *CategoryService) GetAllCategories() ([]*model.ChallengeCategoryModel, error) {
	return s.categoryRepo.GetAllCategories()
}

func (s *CategoryService) UpdateCategory(id uint, dto *dto.UpdateCategoryDTO) (*model.ChallengeCategoryModel, error) {
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	if dto.Name != nil {
		existing, err := s.categoryRepo.GetCategoryByName(*dto.Name)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
		}
		if existing != nil && existing.ID != id {
			return nil, exception.NewConflictException("CATEGORY_ALREADY_EXISTS", "Category", "name", *dto.Name)
		}
		category.Name = *dto.Name
	}
	if dto.Description != nil {
		category.Description = *dto.Description
	}

	return s.categoryRepo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.DeleteCategory(id)
}
