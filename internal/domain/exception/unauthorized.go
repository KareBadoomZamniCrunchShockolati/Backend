package exception

import "net/http"

const (
	ErrorTypeAuthMissingID          = "AUTH_MISSING_ID"
	ErrorTypeAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"

	ErrorTypeAuthHeaderMalformed = "AUTH_HEADER_MALFORMED"
	ErrorTypeTokenEmpty          = "TOKEN_EMPTY"
	ErrorTypeTokenInvalidExpired = "TOKEN_INVALID_EXPIRED"
)

type UnauthorizedException struct {
	*BaseError
}

func (e *UnauthorizedException) ClientError() {}

func NewUnauthorizedException(code string, details map[string]any, params ...string) *UnauthorizedException {
	if details == nil {
		details = map[string]any{}
	}
	return &UnauthorizedException{
		BaseError: NewBaseError(code, http.StatusUnauthorized, details).WithParams(params...),
	}
}

func NewMissingUserIDException() *UnauthorizedException {
	return NewUnauthorizedException(ErrorTypeAuthMissingID, map[string]any{
		"field": "user_id",
	})
}

func NewAuthInvalidCredentials() *UnauthorizedException {
	return NewUnauthorizedException(ErrorTypeAuthInvalidCredentials, nil)
}


func NewAuthHeaderMalformed() *UnauthorizedException {
	return NewUnauthorizedException(ErrorTypeAuthHeaderMalformed, map[string]any{
		"header": "Authorization",
	})
}

func NewTokenEmpty() *UnauthorizedException {
	return NewUnauthorizedException(ErrorTypeTokenEmpty, nil)
}

func NewTokenInvalidOrExpired() *UnauthorizedException {
	return NewUnauthorizedException(ErrorTypeTokenInvalidExpired, nil)
}
