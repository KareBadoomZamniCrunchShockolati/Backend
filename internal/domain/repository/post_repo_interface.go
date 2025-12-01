package repository

import "challenge-app/internal/domain/model"

type PostRepository interface {
	CreatePost(post *model.Post) (*model.Post, error)
	GetPost(postID uint) (*model.Post, error)
	GetPostsByUser(userID uint, offset, limit int) ([]*model.Post, error)
	GetPostsByChallenge(challengeID uint, offset, limit int) ([]*model.Post, error)
	GetFeedPosts(userID uint, offset, limit int) ([]*model.Post, error)
	UpdatePost(post *model.Post) (*model.Post, error)
	DeletePost(postID uint) error
}
