package enum

type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"
	RequestStatusAccepted RequestStatus = "accepted"
	RequestStatusRejected RequestStatus = "rejected"
)

func GetAllRequestStatus() []RequestStatus {
	return []RequestStatus{
		RequestStatusPending,
		RequestStatusAccepted,
		RequestStatusRejected,
	}
}