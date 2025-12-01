package enum

type ParticipantStatus string

const (
	StatusJoined  ParticipantStatus = "joined"
	StatusPending ParticipantStatus = "pending"
	StatusInvited ParticipantStatus = "invited"
	StatusRejected ParticipantStatus = "rejected"
)


func GetAllParticipantStatus() []ParticipantStatus {
	return []ParticipantStatus{
		StatusJoined,
		StatusPending,
		StatusInvited,
	}
}