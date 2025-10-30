package exception

import (
	"net/http"
)

type UnauthorizedException struct {
	BaseError
}

func NewUnauthorizedException(msg string, code string) UnauthorizedException {
	return UnauthorizedException{
		BaseError: BaseError{
			errorCode:  code,
			message:    msg,
			httpStatus: http.StatusUnauthorized, // 401
			details:    nil, 
		},
	}
}
