package exception

import (
	"net/http"
)

type UnauthorizedException struct {
	BaseError
}

func NewUnauthorizedException(msg string, code string) *UnauthorizedException {
	return &UnauthorizedException{
		BaseError: BaseError{
			errorCode:  code,
			message:    msg,
			httpStatus: http.StatusUnauthorized, // 401
			details:    nil, 
		},
	}
}

func NewMissingUserIDException() *UnauthorizedException {
	return NewUnauthorizedException(
		"Authentication failed: User ID not found in context.",
		"AUTH_MISSING_ID",
	)
}

func NewAuthInvalidCredentials() *UnauthorizedException {
	return NewUnauthorizedException(
		"Invalid credentials.",
		"AUTH_INVALID_CREDENTIALS",
	)
}

func (e UnauthorizedException) ClientError() {}