package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type ChallengeInviteRepository struct {
	db *gorm.DB
}

func NewChallengeInviteRepository(db *gorm.DB) *ChallengeInviteRepository {
	return &ChallengeInviteRepository{db: db}
}

func (r *ChallengeInviteRepository) CreateInvite(invite *model.ChallengeInvite) (*model.ChallengeInvite, error) {
	entity := toChallengeInviteEntity(invite)
	result := r.db.Create(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	return toChallengeInviteModelWithID(*entity, invite), nil
}

func (r *ChallengeInviteRepository) GetInvite(inviteID uint) (*model.ChallengeInvite, error) {
	var entity entity.ChallengeInviteEntity
	result := r.db.First(&entity, inviteID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(result.Error)
	}
	return toChallengeInviteModel(&entity), nil
}

func (r *ChallengeInviteRepository) GetUserInvites(userID uint) ([]*model.ChallengeInvite, error) {
	var entities []entity.ChallengeInviteEntity
	result := r.db.Where("invitee_id = ?", userID).Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	invites := make([]*model.ChallengeInvite, len(entities))
	for i, e := range entities {
		invites[i] = toChallengeInviteModel(&e)
	}
	return invites, nil
}

func (r *ChallengeInviteRepository) GetChallengeInvites(challengeID uint) ([]*model.ChallengeInvite, error) {
	var entities []entity.ChallengeInviteEntity
	result := r.db.Where("challenge_id = ?", challengeID).Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	invites := make([]*model.ChallengeInvite, len(entities))
	for i, e := range entities {
		invites[i] = toChallengeInviteModel(&e)
	}
	return invites, nil
}

func (r *ChallengeInviteRepository) UpdateInvite(invite *model.ChallengeInvite) (*model.ChallengeInvite, error) {
	entity := toChallengeInviteEntity(invite)
	entity.ID = invite.ID

	result := r.db.Save(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	invite.UpdatedAt = entity.UpdatedAt
	return invite, nil
}

func (r *ChallengeInviteRepository) DeleteInvite(inviteID uint) error {
	result := r.db.Delete(&entity.ChallengeInviteEntity{}, inviteID)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

// Helpers
func toChallengeInviteEntity(m *model.ChallengeInvite) *entity.ChallengeInviteEntity {
	return &entity.ChallengeInviteEntity{
		ChallengeID: m.ChallengeID,
		InviterID:   m.InviterID,
		InviteeID:   m.InviteeID,
		Status:      uint(m.Status),
	}
}

func toChallengeInviteModel(e *entity.ChallengeInviteEntity) *model.ChallengeInvite {
	return &model.ChallengeInvite{
		ID:          e.ID,
		ChallengeID: e.ChallengeID,
		InviterID:   e.InviterID,
		InviteeID:   e.InviteeID,
		Status:      enum.InviteStatus(e.Status),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toChallengeInviteModelWithID(e entity.ChallengeInviteEntity, original *model.ChallengeInvite) *model.ChallengeInvite {
	original.ID = e.ID
	original.CreatedAt = e.CreatedAt
	original.UpdatedAt = e.UpdatedAt
	return original
}
