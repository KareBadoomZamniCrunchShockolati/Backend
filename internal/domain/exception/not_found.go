// internal/domain/exception/not_found.go
package exception

import (
	"fmt"
	"net/http"
)

type NotFoundException struct {
	*BaseError
}

func NewNotFoundException(item string, id string, code string) *NotFoundException {
	msg := fmt.Sprintf("Resource '%s' with identifier '%s' not found.", item, id)
	return &NotFoundException{
		BaseError: NewBaseError(
			code,
			msg,
			http.StatusNotFound, // 404
			map[string]any{
				"resource_type": item,
				"identifier":    id,
			},
		),
	}
}

func (e NotFoundException) ClientError() {}
