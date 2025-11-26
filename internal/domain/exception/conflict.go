package exception

import (
	"fmt"
	"net/http"
)

type ConflictException struct {
	*BaseError
}

func NewConflictException(resource, field, code string) *ConflictException {
	msg := fmt.Sprintf("The %s already exists for the field %s.", resource, field)
	return &ConflictException{
		BaseError: NewBaseError(
			code,
			msg,
			http.StatusConflict, // 409
			map[string]any{
				"resource": resource,
				"field":    field,
			},
		),
	}
}

func NewUserConflictException(field string) *ConflictException {
	return NewConflictException("User", field, "USER_ALREADY_EXISTS")
}

func (e *ConflictException) ClientError() {}
