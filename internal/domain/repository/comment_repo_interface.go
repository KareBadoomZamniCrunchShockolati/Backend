package repository

import "challenge-app/internal/domain/model"

type CommentRepository interface {
	CreateComment(comment *model.Comment) (*model.Comment, error)
	GetComment(commentID uint) (*model.Comment, error)
	GetComments(entityType model.CommentType, entityID uint, offset, limit int) ([]*model.Comment, error)
	UpdateComment(comment *model.Comment) (*model.Comment, error)
	DeleteComment(commentID uint) error
	GetCommentCount(entityType model.CommentType, entityID uint) (uint, error)
}
