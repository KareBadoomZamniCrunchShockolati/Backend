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

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) postgres.UserRepository {
	return &UserRepository{DB: db}
}

// --- Helper functions for converting between Domain and Infrastructure ---

// toModel converts an Infrastructure entity (with GORM tags) to a Domain model (pure business logic).
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
		CreatedAt:      e.CreatedAt,
		// UpdatedAt را در Domain model قرار نمی‌دهیم، چون جزئیات زیرساخت است.
	}
}

// toEntity converts a Domain model (pure business logic) to an Infrastructure entity (with GORM tags).
func toEntity(m *model.UserModel) *entity.UserEntity {
	if m == nil {
		return nil
	}
	return &entity.UserEntity{
		ID:             m.ID,
		Username:       m.Username,
		Email:          m.Email,
		PasswordHash:   m.PasswordHash,
		Bio:            m.Bio,
		ProfilePicture: m.ProfilePicture,
		CreatedAt:      m.CreatedAt,
		// UpdatedAt توسط GORM به‌طور خودکار پر می‌شود.
	}
}

// --- CRUD Implementation ---

// CreateUser (CRUD - Create) inserts a new user into the database.
func (r *UserRepository) CreateUser(user *model.UserModel)  error{
	// تبدیل مدل دامین به انتیتی زیرساخت
	userEntity := toEntity(user)
	return r.DB.Create(userEntity).Error
	
}

// GetUserByEmail retrieves a user by their email address.
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

// GetAllUsers (CRUD - Read) retrieves all users.
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

// GetUserByID (CRUD - Read) retrieves a user by their ID.
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

// UpdateUser (CRUD - Update) updates an existing user's details.
func (r *UserRepository) UpdateUser(user *model.UserModel) (*model.UserModel, error) {
	userEntity := toEntity(user)
	result := r.DB.Model(&entity.UserEntity{}).Where("id = ?", user.ID).Updates(userEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	return toModel(userEntity), nil
}

// DeleteUser (CRUD - Delete) removes a user from the database by ID.
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
