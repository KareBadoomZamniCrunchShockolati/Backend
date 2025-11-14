package exception

import (
	"net/http"
)

const (
	ErrorTypeAuthMissingID          = "AUTH_MISSING_ID"
	ErrorTypeAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
)

type UnauthorizedException struct {
	*BaseError
}

func NewUnauthorizedException(msg string, code string) *UnauthorizedException {
	return &UnauthorizedException{
		BaseError: NewBaseError(
			code,
			msg,
			http.StatusUnauthorized, // 401
			nil,
		),
	}
}

func NewMissingUserIDException() *UnauthorizedException {
	return NewUnauthorizedException(
		"Authentication failed: User ID not found in context.",
		ErrorTypeAuthMissingID,
	)
}

func NewAuthInvalidCredentials() *UnauthorizedException {
	return NewUnauthorizedException(
		"Invalid credentials.",
		ErrorTypeAuthInvalidCredentials,
	)
}

func (e *UnauthorizedException) ClientError() {}
