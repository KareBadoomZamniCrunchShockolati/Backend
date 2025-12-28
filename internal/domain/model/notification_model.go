package model

import "time"

type NotificationType string

const (
	NotifFollowed             NotificationType = "followed"
	NotifChallengeJoinAccepted NotificationType = "challenge_join_accepted"
	NotifChallengeUserJoined   NotificationType = "challenge_user_joined"
	NotifChallengeReminder     NotificationType = "challenge_reminder"
)

type Notification struct {
	ID        uint
	UserID    uint
	Type      NotificationType
	Title     string
	Body      string
	Data      map[string]any
	CreatedAt time.Time
	ReadAt    *time.Time
}
