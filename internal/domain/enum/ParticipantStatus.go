package enum

type ParticipantStatus uint

const (
	StatusJoined ParticipantStatus = iota + 1
	StatusPending
	StatusInvited
	StatusRejected
)

func (s ParticipantStatus) String() string {
	switch s {
	case StatusJoined:
		return "joined"
	case StatusPending:
		return "pending"
	case StatusInvited:
		return "invited"
	case StatusRejected:
		return "rejected"
	}
	return "unknown"
}

func GetAllParticipantStatuses() []ParticipantStatus {
	return []ParticipantStatus{
		StatusJoined,
		StatusPending,
		StatusInvited,
		StatusRejected,
	}
}
