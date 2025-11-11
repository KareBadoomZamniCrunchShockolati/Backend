package dto

type FollowRequest struct {
	FollowingID uint `json:"following_id" binding:"required"`
}

type FollowResponse struct {
	Message string `json:"message"`
}

type FollowStatusResponse struct {
	IsFollowing bool `json:"is_following"`
}

type FollowStatsResponse struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Count int            `json:"count"`
}

// for when you want to remove one of your followers
type RemoveFollowerRequest struct {
	FollowerID uint `json:"follower_id" binding:"required"`
}

type UserURI struct {
	ID uint `uri:"id" binding:"required"`
}