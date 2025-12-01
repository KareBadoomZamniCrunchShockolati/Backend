package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"encoding/json"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *model.Post) (*model.Post, error) {
	picturesJSON, err := json.Marshal(post.Pictures)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	postEntity := &entity.PostEntity{
		UserID:      post.UserID,
		Description: post.Description,
		ChallengeID: post.ChallengeID,
		Pictures:    string(picturesJSON),
	}

	result := r.db.Create(postEntity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toPostModel(postEntity), nil
}

func (r *PostRepository) GetPost(postID uint) (*model.Post, error) {
	var postEntity entity.PostEntity
	result := r.db.First(&postEntity, postID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toPostModel(&postEntity), nil
}

func (r *PostRepository) GetPostsByUser(userID uint, offset, limit int) ([]*model.Post, error) {
	var entities []entity.PostEntity
	result := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&entities)

	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toPostModels(entities), nil
}

func (r *PostRepository) GetPostsByChallenge(challengeID uint, offset, limit int) ([]*model.Post, error) {
	var entities []entity.PostEntity
	result := r.db.Where("challenge_id = ?", challengeID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&entities)

	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toPostModels(entities), nil
}

func (r *PostRepository) GetFeedPosts(userID uint, offset, limit int) ([]*model.Post, error) {
	var entities []entity.PostEntity
	result := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&entities)

	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toPostModels(entities), nil
}

func (r *PostRepository) UpdatePost(post *model.Post) (*model.Post, error) {
	var postEntity entity.PostEntity
	result := r.db.First(&postEntity, post.ID)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	picturesJSON, err := json.Marshal(post.Pictures)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	postEntity.Description = post.Description
	postEntity.Pictures = string(picturesJSON)

	result = r.db.Save(&postEntity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toPostModel(&postEntity), nil
}

func (r *PostRepository) DeletePost(postID uint) error {
	result := r.db.Delete(&entity.PostEntity{}, postID)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

// Helper functions
func toPostModel(e *entity.PostEntity) *model.Post {
	var pictures []string
	if e.Pictures != "" {
		json.Unmarshal([]byte(e.Pictures), &pictures)
	}

	return &model.Post{
		ID:          e.ID,
		UserID:      e.UserID,
		Description: e.Description,
		ChallengeID: e.ChallengeID,
		Pictures:    pictures,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toPostModels(entities []entity.PostEntity) []*model.Post {
	posts := make([]*model.Post, len(entities))
	for i, e := range entities {
		posts[i] = toPostModel(&e)
	}
	return posts
}
