package postgres

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository/postgres"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"errors"

	"gorm.io/gorm"
)

// UserRepository implements the domain.UserRepository interface.
type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) postgres.UserRepository {
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
	}
}

// --- CRUD Implementation ---

func (r *UserRepository) CreateUser(user *model.UserModel) error {
    userEntity := toEntity(user)
    if err := r.DB.Create(userEntity).Error; err != nil {
        return err
    }
    user.ID = userEntity.ID 
    return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.UserModel, error) {
	var userEntity entity.UserEntity
	result := r.DB.Where("email = ?", email).First(&userEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return toModel(&userEntity), nil
}

func (r *UserRepository) GetAllUsers() ([]model.UserModel, error) {
	var userEntities []entity.UserEntity
	result := r.DB.Find(&userEntities)
	if result.Error != nil {
		return nil, result.Error
	}
	var userModels []model.UserModel
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
			return nil, result.Error
		}
		return nil, result.Error
	}
	return toModel(&userEntity), nil
}

func (r *UserRepository) UpdateUser(user *model.UserModel) (*model.UserModel, error) {
	userEntity := toEntity(user)
	result := r.DB.Model(&entity.UserEntity{}).Where("id = ?", user.ID).Updates(userEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	return toModel(userEntity), nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	result := r.DB.Delete(&entity.UserEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}
	return nil
}
