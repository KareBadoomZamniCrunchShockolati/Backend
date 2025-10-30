package exception

import (
	"fmt"
	"net/http"
)

type ConflictException struct {
	BaseError
}

func NewConflictException(resource string, field string, code string) ConflictException {
	msg := fmt.Sprintf("The %s already exists for the field %s.", resource, field)
	return ConflictException{
		BaseError: BaseError{
			errorCode:  code,
			message:    msg,
			httpStatus: http.StatusConflict, // 409
			details: map[string]any{
				"resource": resource,
				"field":    field,
			},
		},
	}
}
