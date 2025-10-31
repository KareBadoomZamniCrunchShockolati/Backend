package repository

import "challenge-app/internal/domain/model"

type FollowRepository interface {
	Follow(followerID, followingID uint) error
	Unfollow(followerID, followingID uint) error
	RemoveFollower(userID, followerID uint) error
	GetFollowers(userID uint) ([]model.UserModel, error)
	GetFollowing(userID uint) ([]model.UserModel, error)
	IsFollowing(followerID, followingID uint) (bool, error)
	GetFollowersCount(userID uint) (int64, error)
	GetFollowingCount(userID uint) (int64, error)
}
