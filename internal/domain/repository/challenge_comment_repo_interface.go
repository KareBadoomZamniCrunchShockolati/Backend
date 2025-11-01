package repository

import "challenge-app/internal/domain/model"

type ChallengeCommentRepository interface {
	CreateComment(comment *model.ChallengeComment) error
	UpdateComment(comment *model.ChallengeComment) error
	DeleteComment(commentID uint) error
	ListComments(challengeID uint, offset, limit int) ([]*model.ChallengeComment, error)
}
