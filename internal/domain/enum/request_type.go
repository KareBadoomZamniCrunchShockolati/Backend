package enum

type RequestType string

const (
	RequestTypeInvite  RequestType = "invite"
	RequestTypeRequest RequestType = "request"
)

func GetAllRequestType() []RequestType {
	return []RequestType{
		RequestTypeInvite,
		RequestTypeRequest,
	}
}