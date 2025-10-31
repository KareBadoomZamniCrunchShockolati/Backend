package serviceinterface

import "challenge-app/internal/domain/model"

type FollowServicer interface {
	Follow(followerID, followingID uint) error
	Unfollow(followerID, followingID uint) error
	RemoveFollower(userID, followerID uint) error
	GetFollowers(userID uint) ([]model.UserModel, error)
	GetFollowing(userID uint) ([]model.UserModel, error)
	IsFollowing(followerID, followingID uint) (bool, error)
	GetFollowStats(userID uint) (*FollowStats, error)
}

type FollowStats struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}
