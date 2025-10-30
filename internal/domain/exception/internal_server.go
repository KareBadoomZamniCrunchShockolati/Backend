package exception

import "net/http"

type InternalServerException struct {
	BaseError
}

func NewInternalServerException(clientMsg string, code string) InternalServerException {
	return InternalServerException{
		BaseError: BaseError{
			errorCode:  code,
			message:    clientMsg,
			httpStatus: http.StatusInternalServerError, // 500
			details:    nil,
		},
	}
}
