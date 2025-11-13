package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
)

type CategoryService struct {
	Repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{Repo: repo}
}

func (s *CategoryService) CreateCategory(dto *dto.CreateCategoryDTO) (*model.ChallengeCategoryModel, error) {
	existing, err := s.Repo.GetCategoryByName(dto.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, exception.NewConflictException("Category", "name", "CATEGORY_ALREADY_EXISTS")
	}

	return s.Repo.CreateCategory(&model.ChallengeCategoryModel{
		Name:        dto.Name,
		Description: dto.Description,
	})
}

func (s *CategoryService) GetCategoryByID(id uint) (*model.ChallengeCategoryModel, error) {
	return s.Repo.GetCategoryByID(id)
}

func (s *CategoryService) GetAllCategories() ([]*model.ChallengeCategoryModel, error) {
	return s.Repo.GetAllCategories()
}

func (s *CategoryService) UpdateCategory(id uint, dto *dto.UpdateCategoryDTO) (*model.ChallengeCategoryModel, error) {
	category, err := s.Repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		existing, _ := s.Repo.GetCategoryByName(*dto.Name)
		if existing != nil && existing.ID != id {
			return nil, exception.NewConflictException("Category", "name", "CATEGORY_ALREADY_EXISTS")
		}
		category.Name = *dto.Name
	}
	if dto.Description != nil {
		category.Description = *dto.Description
	}

	return s.Repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(id uint) error {
	return s.Repo.DeleteCategory(id)
}
