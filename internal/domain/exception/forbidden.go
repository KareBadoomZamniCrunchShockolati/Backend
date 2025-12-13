// internal/domain/exception/forbidden.go
package exception

import (
	"net/http"
)

type ForbiddenException struct {
	*BaseError
}

func NewForbiddenException(msg string, code string) *ForbiddenException {
	return &ForbiddenException{
		BaseError: NewBaseError(
			code,
			msg,
			http.StatusForbidden, // 403
			nil,
		),
	}
}

func (e *ForbiddenException) ClientError() {}
