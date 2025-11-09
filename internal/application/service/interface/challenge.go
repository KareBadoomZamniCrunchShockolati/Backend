package serviceinterface

import (
    "challenge-app/internal/application/dto"
    "challenge-app/internal/domain/model"
)

type ChallengeServicer interface {
    // Challenge CRUD
    CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error)
    UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error)
    GetChallengeByID(id uint) (*model.ChallengeModel, error)
    DeleteChallenge(challengeID uint, currentUserID uint) error
    // StopChallenge(challengeID uint) error

    // Listing challenges
    // ListPublicChallenges(offset, limit int) ([]*model.ChallengeModel, error)
    // ListPrivateChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
    // ListInviteChallengesForUser(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
    // ListByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
    // ListByCategory(category string, offset, limit int) ([]*model.ChallengeModel, error)

    // Participants
    // SubscribeChallenge(userID, challengeID uint) error
    // UnsubscribeChallenge(userID, challengeID uint) error
    // ListChallengeParticipants(challengeID uint) ([]*model.ChallengeParticipant, error)
    // ListParticipantsInUserFollowings(challengeID, userID uint) ([]*model.ChallengeParticipant, error)

    // Comments
    // AddComment(comment *model.ChallengeComment) error
    // UpdateComment(comment *model.ChallengeComment) error
    // DeleteComment(commentID uint) error
    // ListComments(challengeID uint, offset, limit int) ([]*model.ChallengeComment, error)
}
