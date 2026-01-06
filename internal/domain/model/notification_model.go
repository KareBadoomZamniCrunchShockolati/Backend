package model

import "time"

type NotificationType string

const (
	NotifFollowed                   NotificationType = "followed"
	NotifJoinRequestSent            NotificationType = "challenge_join_request_sent"
	NotifJoinRequestAccepted        NotificationType = "challenge_join_request_accepted"
	NotifInviteSent                 NotificationType = "challenge_invite_sent"
	NotifInviteAccepted             NotificationType = "challenge_invite_accepted"
	NotifChallengeReminder          NotificationType = "challenge_reminder"
	NotifChallengeParticipantJoined NotificationType = "challenge_participant_joined"
	NotifChallengeCommented         NotificationType = "challenge_commented"
	NotifChallengeLiked             NotificationType = "challenge_liked"
)

type Notification struct {
	ID     uint
	UserID uint
	Type   NotificationType
	TitleKey string
	BodyKey  string
	Data     map[string]any
	CreatedAt time.Time
	ReadAt    *time.Time
}
