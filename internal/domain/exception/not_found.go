// internal/domain/exception/not_found.go
package exception

import (
	"fmt"
	"net/http"
)

type NotFoundException struct {
	BaseError
}

func NewNotFoundException(item string, id string, code string) NotFoundException {
	msg := fmt.Sprintf("Resource '%s' with identifier '%s' not found.", item, id)
	return NotFoundException{
		BaseError: BaseError{
			errorCode:  code,
			message:    msg,
			httpStatus: http.StatusNotFound, // 404
			details: map[string]any{
				"resource_type": item,
				"identifier":    id,
			},
		},
	}
}

func (e NotFoundException) ClientError() {}
