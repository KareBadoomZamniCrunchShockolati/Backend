package service

import (
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"fmt"
)

type FollowService struct {
	FollowRepo repository.FollowRepository
	UserRepo   repository.UserRepository
}

func NewFollowService(followRepo repository.FollowRepository, userRepo repository.UserRepository) *FollowService {
	return &FollowService{
		FollowRepo: followRepo,
		UserRepo:   userRepo,
	}
}

func (s *FollowService) Follow(followerID, followingID uint) error {
	// Check if target user exists
	_, err := s.UserRepo.GetUserByID(followingID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	return s.FollowRepo.Follow(followerID, followingID)
}

func (s *FollowService) Unfollow(followerID, followingID uint) error {
	return s.FollowRepo.Unfollow(followerID, followingID)
}

func (s *FollowService) GetFollowers(userID uint) ([]model.UserModel, error) {
	return s.FollowRepo.GetFollowers(userID)
}

func (s *FollowService) GetFollowing(userID uint) ([]model.UserModel, error) {
	return s.FollowRepo.GetFollowing(userID)
}

func (s *FollowService) IsFollowing(followerID, followingID uint) (bool, error) {
	return s.FollowRepo.IsFollowing(followerID, followingID)
}

func (s *FollowService) GetFollowStats(userID uint) (*serviceinterface.FollowStats, error) {
	followers, err := s.FollowRepo.GetFollowersCount(userID)
	if err != nil {
		return nil, err
	}

	following, err := s.FollowRepo.GetFollowingCount(userID)
	if err != nil {
		return nil, err
	}

	return &serviceinterface.FollowStats{
		FollowersCount: followers,
		FollowingCount: following,
	}, nil
}

func (s *FollowService) RemoveFollower(userID, followerID uint) error {
	// Check if the follower actually exists
	_, err := s.UserRepo.GetUserByID(followerID)
	if err != nil {
		return fmt.Errorf("follower user not found")
	}

	// Check if the follower is actually following the user
	isFollowing, err := s.FollowRepo.IsFollowing(followerID, userID)
	if err != nil {
		return fmt.Errorf("error checking follow relationship: %w", err)
	}

	if !isFollowing {
		return fmt.Errorf("user is not following you")
	}

	// Remove the follower
	return s.FollowRepo.RemoveFollower(userID, followerID)
}
