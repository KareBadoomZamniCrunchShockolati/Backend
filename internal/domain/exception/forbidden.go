package exception

import "net/http"

const (
	ErrorTypeForbidden = "FORBIDDEN"
)

type ForbiddenException struct {
	*BaseError
}

func NewForbiddenException(code string, params ...string) *ForbiddenException {
	return &ForbiddenException{
		BaseError: NewBaseError(code, http.StatusForbidden, map[string]any{}).
			WithParams(params...),
	}
}

func (e *ForbiddenException) ClientError() {}
