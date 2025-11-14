package repository

import "challenge-app/internal/domain/model"

type ChallengeCommentRepository interface {
	CreateComment(comment *model.ChallengeComment) (*model.ChallengeComment, error)
	GetComment(commentID uint) (*model.ChallengeComment, error)
	GetChallengeComments(challengeID uint) ([]*model.ChallengeComment, error)
	UpdateComment(comment *model.ChallengeComment) (*model.ChallengeComment, error)
	DeleteComment(commentID uint) error
}
