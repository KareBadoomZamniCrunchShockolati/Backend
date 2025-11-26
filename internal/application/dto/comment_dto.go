package dto

import "time"

type AddCommentDTO struct {
	ChallengeID uint   `json:"challenge_id" validate:"required"`
	Content     string `json:"content" validate:"required,min=1,max=1000"`
}

type CommentResponseDTO struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:user_id`
	ChallengeID uint      `json:"challenge_id"`
	Username    string    `json:"username"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}
