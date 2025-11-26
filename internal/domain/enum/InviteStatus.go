package enum

type InviteStatus uint

const (
	InvitePending InviteStatus = iota + 1
	InviteAccepted
	InviteDeclined
	InviteExpired
)

func (i InviteStatus) String() string {
	switch i {
	case InvitePending:
		return "pending"
	case InviteAccepted:
		return "accepted"
	case InviteDeclined:
		return "declined"
	case InviteExpired:
		return "expired"
	}
	return "unknown"
}

func GetAllInviteStatuses() []InviteStatus {
	return []InviteStatus{
		InvitePending,
		InviteAccepted,
		InviteDeclined,
		InviteExpired,
	}
}
