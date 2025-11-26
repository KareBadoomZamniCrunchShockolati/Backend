package enum

type RequestStatus uint

const (
	RequestPending RequestStatus = iota + 1
	RequestAccepted
	RequestRejected
)

func (r RequestStatus) String() string {
	switch r {
	case RequestPending:
		return "pending"
	case RequestAccepted:
		return "accepted"
	case RequestRejected:
		return "rejected"
	}
	return "unknown"
}

func GetAllRequestStatuses() []RequestStatus {
	return []RequestStatus{
		RequestPending,
		RequestAccepted,
		RequestRejected,
	}
}
