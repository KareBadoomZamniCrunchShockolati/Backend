package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(comment *model.Comment) (*model.Comment, error) {
	commentEntity := &entity.CommentEntity{
		EntityType: string(comment.EntityType),
		EntityID:   comment.EntityID,
		UserID:     comment.UserID,
		Content:    comment.Content,
		ParentID:   comment.ParentID,
	}

	result := r.db.Create(commentEntity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toCommentModel(commentEntity), nil
}

func (r *CommentRepository) GetComment(commentID uint) (*model.Comment, error) {
	var commentEntity entity.CommentEntity
	result := r.db.First(&commentEntity, commentID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toCommentModel(&commentEntity), nil
}

func (r *CommentRepository) GetComments(entityType model.CommentType, entityID uint, offset, limit int) ([]*model.Comment, error) {
	var entities []entity.CommentEntity
	result := r.db.Where("entity_type = ? AND entity_id = ?", string(entityType), entityID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&entities)

	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toCommentModels(entities), nil
}

func (r *CommentRepository) UpdateComment(comment *model.Comment) (*model.Comment, error) {
	var commentEntity entity.CommentEntity
	result := r.db.First(&commentEntity, comment.ID)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	commentEntity.Content = comment.Content

	result = r.db.Save(&commentEntity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	return toCommentModel(&commentEntity), nil
}

func (r *CommentRepository) DeleteComment(commentID uint) error {
	result := r.db.Delete(&entity.CommentEntity{}, commentID)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

func (r *CommentRepository) GetCommentCount(entityType model.CommentType, entityID uint) (uint, error) {
	var count int64
	result := r.db.Model(&entity.CommentEntity{}).
		Where("entity_type = ? AND entity_id = ?", string(entityType), entityID).
		Count(&count)
	if result.Error != nil {
		return 0, exception.NewRepositoryError(result.Error)
	}
	return uint(count), nil
}

// Helper functions
func toCommentModel(e *entity.CommentEntity) *model.Comment {
	return &model.Comment{
		ID:         e.ID,
		EntityType: model.CommentType(e.EntityType),
		EntityID:   e.EntityID,
		UserID:     e.UserID,
		Content:    e.Content,
		ParentID:   e.ParentID,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func toCommentModels(entities []entity.CommentEntity) []*model.Comment {
	comments := make([]*model.Comment, len(entities))
	for i, e := range entities {
		comments[i] = toCommentModel(&e)
	}
	return comments
}
