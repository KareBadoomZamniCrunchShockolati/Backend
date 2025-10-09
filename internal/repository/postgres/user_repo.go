package postgres

import (
	"challenge-app/internal/domain"
	"challenge-app/internal/repository/models"
    "github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthRepository implements the domain.UserRepository interface.
type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{DB: db}
}

// Helper functions for converting between Domain and Model (essential for decoupling)
func toModel(u *domain.User) models.UserModel {
    return models.UserModel{
        ID: u.ID, Username: u.Username, Email: u.Email,
        PasswordHash: u.PasswordHash, Bio: u.Bio, CreatedAt: u.CreatedAt,
    }
}
func toDomain(u *models.UserModel) *domain.User {
    return &domain.User{
        ID: u.ID, Username: u.Username, Email: u.Email,
        PasswordHash: u.PasswordHash, Bio: u.Bio, CreatedAt: u.CreatedAt,
    }
}

// --- CRUD Implementation ---

// CreateUser (CRUD - Create)
func (r *UserRepository) CreateUser(user *domain.User) error {
	model := toModel(user)
	return r.DB.Create(&model).Error
}

// GetUserByEmail 
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var model models.UserModel
	result := r.DB.Where("email = ?", email).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(&model), nil
}

func (r *UserRepository) GetAllUsers() ([]domain.User, error) {
	var models []models.UserModel
	if err := r.DB.Find(&models).Error; err != nil {
		return nil, err
	}

	var users []domain.User
	for _, model := range models {
		users = append(users, *toDomain(&model))
	}
	return users, nil
}

// GetUserByID (CRUD - Read)
func (r *UserRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
    var model models.UserModel
    result := r.DB.First(&model, "id = ?", id)
    if result.Error != nil {
        return nil, result.Error
    }
    return toDomain(&model), nil
}

// UpdateUser (CRUD - Update)
func (r *UserRepository) UpdateUser(user *domain.User) error {
    

    updates := map[string]interface{}{
        "username": user.Username,
        "email":    user.Email, 
        "bio":      user.Bio,
    }

    result := r.DB.Model(&models.UserModel{}).
        Where("id = ?", user.ID).
        Updates(updates) 
    
    if result.Error != nil {
        return result.Error
    }
    
    // Check if the record was actually modified
    if result.RowsAffected == 0 {
        var count int64
        r.DB.Model(&models.UserModel{}).Where("id = ?", user.ID).Count(&count)
        if count == 0 {
            return gorm.ErrRecordNotFound
        }
    }

    return nil
}
// DeleteUser (CRUD - Delete)
func (r *UserRepository) DeleteUser(id uuid.UUID) error {
    return r.DB.Delete(&models.UserModel{}, "id = ?", id).Error
}