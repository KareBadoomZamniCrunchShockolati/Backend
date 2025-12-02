package postgres

import (
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type ChallengeParticipationRepository struct {
	db *gorm.DB
}

func NewChallengeParticipationRepository(db *gorm.DB) *ChallengeParticipationRepository {
	return &ChallengeParticipationRepository{db: db}
}

func (r *ChallengeParticipationRepository) CreateInvite(inviterID, challengeID, inviteeID uint) (*model.ChallengeInvite, error) {
	ent := &entity.ChallengeParticipationRequestEntity{
		ChallengeID: challengeID,
		FromUserID:  inviterID,
		ToUserID:    inviteeID,
		Type:        enum.RequestTypeInvite,
		Status:      enum.RequestStatusPending,
	}
	if err := r.db.Create(ent).Error; err != nil {
		return nil, err
	}
	return &model.ChallengeInvite{
		ID:          ent.ID,
		ChallengeID: ent.ChallengeID,
		InviterID:   ent.FromUserID,
		InviteeID:   ent.ToUserID,
		Status:      ent.Status,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, nil
}

func (r *ChallengeParticipationRepository) GetInvite(inviteID uint) (*model.ChallengeInvite, error) {
	var ent entity.ChallengeParticipationRequestEntity
	if err := r.db.
		Where("id = ? AND type = ?", inviteID, enum.RequestTypeInvite).
		First(&ent).Error; err != nil {
		return nil, err
	}
	return &model.ChallengeInvite{
		ID:          ent.ID,
		ChallengeID: ent.ChallengeID,
		InviterID:   ent.FromUserID,
		InviteeID:   ent.ToUserID,
		Status:      ent.Status,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, nil
}

func (r *ChallengeParticipationRepository) GetInvitesSentToUser(userID uint) ([]*model.ChallengeInvite, error) {
	var ents []entity.ChallengeParticipationRequestEntity
	if err := r.db.
		Where("to_user_id = ? AND type = ?", userID, enum.RequestTypeInvite).
		Find(&ents).Error; err != nil {
		return nil, err
	}
	invites := make([]*model.ChallengeInvite, len(ents))
	for i := range ents {
		e := &ents[i]
		invites[i] = &model.ChallengeInvite{
			ID:          e.ID,
			ChallengeID: e.ChallengeID,
			InviterID:   e.FromUserID,
			InviteeID:   e.ToUserID,
			Status:      e.Status,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		}
	}
	return invites, nil
}

func (r *ChallengeParticipationRepository) GetInvitesSentFromChallenge(challengeID uint) ([]*model.ChallengeInvite, error) {
	var ents []entity.ChallengeParticipationRequestEntity
	if err := r.db.
		Where("challenge_id = ? AND type = ?", challengeID, enum.RequestTypeInvite).
		Find(&ents).Error; err != nil {
		return nil, err
	}
	invites := make([]*model.ChallengeInvite, len(ents))
	for i := range ents {
		e := &ents[i]
		invites[i] = &model.ChallengeInvite{
			ID:          e.ID,
			ChallengeID: e.ChallengeID,
			InviterID:   e.FromUserID,
			InviteeID:   e.ToUserID,
			Status:      e.Status,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		}
	}
	return invites, nil
}

func (r *ChallengeParticipationRepository) UpdateInvite(invite *model.ChallengeInvite) (*model.ChallengeInvite, error) {
	ent := &entity.ChallengeParticipationRequestEntity{
		ChallengeID: invite.ChallengeID,
		FromUserID:  invite.InviterID,
		ToUserID:    invite.InviteeID,
		Type:        enum.RequestTypeInvite,
		Status:      invite.Status,
	}
	if err := r.db.Save(ent).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (r *ChallengeParticipationRepository) DeleteInvite(inviteID uint) error {
	return r.db.
		Where("id = ? AND type = ?", inviteID, enum.RequestTypeInvite).
		Delete(&entity.ChallengeParticipationRequestEntity{}).Error
}

func (r *ChallengeParticipationRepository) CreateRequest(requesterID, challengeID uint) (*model.ChallengeRequest, error) {
	var challenge entity.ChallengeEntity
	if err := r.db.Where("id = ?", challengeID).First(&challenge).Error; err != nil {
		return nil, err
	}
	ent := &entity.ChallengeParticipationRequestEntity{
		ChallengeID: challengeID,
		FromUserID:  requesterID,
		ToUserID:    challenge.CreatorID, 
		Type:        enum.RequestTypeRequest,
		Status:      enum.RequestStatusPending,
	}
	if err := r.db.Create(ent).Error; err != nil {
		return nil, err
	}
	return &model.ChallengeRequest{
		ID:           ent.ID,
		ChallengeID:  ent.ChallengeID,
		RequesterID:  ent.FromUserID,
		OwnerID:      ent.ToUserID,
		Status:       ent.Status,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
	}, nil
}

func (r *ChallengeParticipationRepository) GetRequest(requestID uint) (*model.ChallengeRequest, error) {
	var ent entity.ChallengeParticipationRequestEntity
	if err := r.db.
		Where("id = ? AND type = ?", requestID, enum.RequestTypeRequest).
		First(&ent).Error; err != nil {
		return nil, err
	}
	return &model.ChallengeRequest{
		ID:           ent.ID,
		ChallengeID:  ent.ChallengeID,
		RequesterID:  ent.FromUserID,
		OwnerID:      ent.ToUserID,
		Status:       ent.Status,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
	}, nil
}

func (r *ChallengeParticipationRepository) GetRequestsSentByUser(userID uint) ([]*model.ChallengeRequest, error) {
	var ents []entity.ChallengeParticipationRequestEntity
	if err := r.db.
		Where("from_user_id = ? AND type = ?", userID, enum.RequestTypeRequest).
		Find(&ents).Error; err != nil {
		return nil, err
	}
	requests := make([]*model.ChallengeRequest, len(ents))
	for i := range ents {
		e := &ents[i]
		requests[i] = &model.ChallengeRequest{
			ID:           e.ID,
			ChallengeID:  e.ChallengeID,
			RequesterID:  e.FromUserID,
			OwnerID:      e.ToUserID,
			Status:       e.Status,
			CreatedAt:    e.CreatedAt,
			UpdatedAt:    e.UpdatedAt,
		}
	}
	return requests, nil
}

func (r *ChallengeParticipationRepository) GetRequestsSentToChallenge(challengeID uint) ([]*model.ChallengeRequest, error) {
	var ents []entity.ChallengeParticipationRequestEntity
	if err := r.db.
		Where("challenge_id = ? AND type = ?", challengeID, enum.RequestTypeRequest).
		Find(&ents).Error; err != nil {
		return nil, err
	}
	requests := make([]*model.ChallengeRequest, len(ents))
	for i := range ents {
		e := &ents[i]
		requests[i] = &model.ChallengeRequest{
			ID:           e.ID,
			ChallengeID:  e.ChallengeID,
			RequesterID:  e.FromUserID,
			OwnerID:      e.ToUserID,
			Status:       e.Status,
			CreatedAt:    e.CreatedAt,
			UpdatedAt:    e.UpdatedAt,
		}
	}
	return requests, nil
}

func (r *ChallengeParticipationRepository) UpdateRequest(request *model.ChallengeRequest) (*model.ChallengeRequest, error) {
	ent := &entity.ChallengeParticipationRequestEntity{
		ChallengeID:  request.ChallengeID,
		FromUserID:   request.RequesterID,
		ToUserID:     request.OwnerID,
		Type:         enum.RequestTypeRequest,
		Status:       request.Status,
	}
	if err := r.db.Save(ent).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r *ChallengeParticipationRepository) DeleteRequest(requestID uint) error {
	return r.db.
		Where("id = ? AND type = ?", requestID, enum.RequestTypeRequest).
		Delete(&entity.ChallengeParticipationRequestEntity{}).Error
}