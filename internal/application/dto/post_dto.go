package dto

import "time"

type CreatePostDTO struct {
	Description string   `json:"description" validate:"required,min=1,max=1000"`
	ChallengeID *uint    `json:"challenge_id,omitempty"`
	Pictures    []string `json:"pictures,omitempty" validate:"max=10"`
	TempKeys    []string `json:"temp_keys,omitempty"`
}

type UpdatePostDTO struct {
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1,max=1000"`
	Pictures    *[]string `json:"pictures,omitempty" validate:"omitempty,max=5"`
}

type PostResponseDTO struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Description  string    `json:"description"`
	ChallengeID  *uint     `json:"challenge_id,omitempty"`
	Pictures     []string  `json:"pictures"`
	LikeCount    uint      `json:"like_count"`
	CommentCount uint      `json:"comment_count"`
	IsLiked      bool      `json:"is_liked"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
