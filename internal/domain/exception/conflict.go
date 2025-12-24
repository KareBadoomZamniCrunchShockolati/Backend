package exception

import "net/http"

const (
	ErrorUserAlreadyExists         = "USER_ALREADY_EXISTS"         
	ErrorUserAlreadyExistsEmail    = "USER_ALREADY_EXISTS_EMAIL"
	ErrorUserAlreadyExistsUsername = "USER_ALREADY_EXISTS_USERNAME"
)

type ConflictException struct {
	*BaseError
}

func (e *ConflictException) ClientError() {}

func NewConflictException(code, resource, field, identifier string) *ConflictException {
	details := map[string]any{
		"resource":   resource,
		"field":      field,
		"identifier": identifier,
	}

	return &ConflictException{
		BaseError: NewBaseError(code, http.StatusConflict, details).
			WithParams(resource, field, identifier),
	}
}

func NewUserAlreadyExistsByEmail(email string) *ConflictException {
	return NewConflictException(ErrorUserAlreadyExists, "User", "email", email)
}

func NewUserAlreadyExistsByUsername(username string) *ConflictException {
	return NewConflictException(ErrorUserAlreadyExists, "User", "username", username)
}
