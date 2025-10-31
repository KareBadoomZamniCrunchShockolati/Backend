package postgres

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository implements the domain.UserRepository interface.
type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// --- Helper functions for converting between Domain and Infrastructure ---

func toModel(e *entity.UserEntity) *model.UserModel {
	if e == nil {
		return nil
	}
	return &model.UserModel{
		ID:             e.ID,
		Username:       e.Username,
		Email:          e.Email,
		PasswordHash:   e.PasswordHash,
		Bio:            e.Bio,
		ProfilePicture: e.ProfilePicture,
		Verified:       e.Verified,
	}
}

func toEntity(m *model.UserModel) *entity.UserEntity {
	if m == nil {
		return nil
	}
	return &entity.UserEntity{
		Username:       m.Username,
		Email:          m.Email,
		PasswordHash:   m.PasswordHash,
		Bio:            m.Bio,
		ProfilePicture: m.ProfilePicture,
		Verified:       m.Verified,
	}
}

// --- CRUD Implementation ---

func (r *UserRepository) CreateUser(user *model.UserModel) error {
	userEntity := toEntity(user)
	if err := r.DB.Create(userEntity).Error; err != nil {
		return fmt.Errorf("database error creating user: %w", err)
	}
	user.ID = userEntity.ID
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.UserModel, error) {
	var userEntity entity.UserEntity
	result := r.DB.Where("email = ?", email).First(&userEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("database error fetching user by email: %w", result.Error)
	}
	return toModel(&userEntity), nil
}

func (r *UserRepository) GetAllUsers() ([]model.UserModel, error) {
	var userEntities []entity.UserEntity
	result := r.DB.Find(&userEntities)
	if result.Error != nil {
		return nil, fmt.Errorf("database error fetching all users: %w", result.Error)
	}
	var userModels []model.UserModel
	if len(userEntities) == 0 {
		return []model.UserModel{}, nil
	}
	for _, e := range userEntities {
		userModels = append(userModels, *toModel(&e))
	}
	return userModels, nil
}

func (r *UserRepository) GetUserByID(id uint) (*model.UserModel, error) {
	var userEntity entity.UserEntity
	result := r.DB.First(&userEntity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("database error fetching user by ID: %w", result.Error)
	}
	return toModel(&userEntity), nil
}

func (r *UserRepository) UpdateUser(user *model.UserModel) (*model.UserModel, error) {
	userEntity := toEntity(user)
	result := r.DB.Model(&entity.UserEntity{}).Where("id = ?", user.ID).Updates(userEntity)
	if result.Error != nil {
		return nil, fmt.Errorf("database error updating user: %w", result.Error)
	}
	return toModel(userEntity), nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	var userEntity entity.UserEntity
    result := r.DB.Unscoped().First(&userEntity, id)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil 
    }
    
    if result.Error != nil {
        return fmt.Errorf("database error fetching user before hard delete: %w", result.Error)
    }

    deleteResult := r.DB.Unscoped().Delete(&entity.UserEntity{}, id)
    
    if deleteResult.Error != nil {
        return fmt.Errorf("database error performing hard delete: %w", deleteResult.Error)
    }
	return nil
}
