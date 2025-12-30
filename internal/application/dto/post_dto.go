package dto

import "time"

type PostImageDTO struct {
	TempKey string `json:"temp_key" validate:"required"`
	URL     string `json:"url,omitempty"`
}

type CreatePostDTO struct {
	Description string         `json:"description"`
	Images      []PostImageDTO `json:"images,omitempty" validate:"max=10"`
	ChallengeID *uint          `json:"challenge_id,omitempty"`
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
