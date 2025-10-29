package postgres

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type FollowRepository struct {
	DB *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	follow := &entity.FollowEntity{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	result := r.DB.Create(follow)
	if result.Error != nil {
		// Check if user is followed before
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("already following this user")
		}
		return fmt.Errorf("database error while creating follow: %v", result.Error)
	}

	return nil
}

func (r *FollowRepository) Unfollow(followerID, followingID uint) error {
	result := r.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entity.FollowEntity{})
	if result.Error != nil {
		return fmt.Errorf("database error while unfollowing: %v", result.Error)
	}
	//check if the user was actually followed before
	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}
	return nil
}
func (r *FollowRepository) RemoveFollower(userID, followerID uint) error {

	result := r.DB.Where("follower_id = ? AND following_id = ?", followerID, userID).Delete(&entity.FollowEntity{})
	if result.Error != nil {
		return fmt.Errorf("database error while removing follower: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("follower relationship not found")
	}
	return nil
}

func (r *FollowRepository) GetFollowers(userID uint) ([]model.UserModel, error) {
	var userEntities []entity.UserEntity

	err := r.DB.
		Table("users").
		Joins("INNER JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ? AND follows.deleted_at IS NULL", userID).
		Find(&userEntities).Error

	if err != nil {
		return nil, fmt.Errorf("database error while fetching followers: %v", err)
	}

	return convertUserEntitiesToModels(userEntities), nil
}

func (r *FollowRepository) GetFollowing(userID uint) ([]model.UserModel, error) {
	var userEntities []entity.UserEntity

	err := r.DB.
		Table("users").
		Joins("INNER JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ? AND follows.deleted_at IS NULL", userID).
		Find(&userEntities).Error

	if err != nil {
		return nil, fmt.Errorf("database error while fetching following: %v", err)
	}

	return convertUserEntitiesToModels(userEntities), nil
}

func (r *FollowRepository) IsFollowing(followerID, followingID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.FollowEntity{}).
		Where("follower_id = ? AND following_id = ? AND deleted_at IS NULL", followerID, followingID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("database error while checking follow: %v", err)
	}

	return count > 0, nil
}

// two functions to be shown in the user profile
func (r *FollowRepository) GetFollowersCount(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&entity.FollowEntity{}).
		Where("following_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *FollowRepository) GetFollowingCount(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&entity.FollowEntity{}).
		Where("follower_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func convertUserEntitiesToModels(entities []entity.UserEntity) []model.UserModel {
	var models []model.UserModel
	for _, e := range entities {
		models = append(models, model.UserModel{
			ID:             e.ID,
			Username:       e.Username,
			Email:          e.Email,
			Bio:            e.Bio,
			ProfilePicture: e.ProfilePicture,
			Verified:       e.Verified,
		})
	}
	return models
}
