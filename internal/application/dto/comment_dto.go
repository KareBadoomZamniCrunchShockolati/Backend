package dto

type AddCommentDTO struct {
	ChallengeID uint   `json:"challenge_id" validate:"required"`
	Content     string `json:"content" validate:"required,min=1,max=1000"`
}
