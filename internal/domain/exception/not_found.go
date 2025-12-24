package exception

import "net/http"

type NotFoundException struct {
	*BaseError
}

func (e *NotFoundException) ClientError() {}

func NewNotFoundException(resource string, identifier string, code string, params ...string) *NotFoundException {
	details := map[string]any{
		"resource":   resource,
		"identifier": identifier,
	}

	return &NotFoundException{
		BaseError: NewBaseError(code, http.StatusNotFound, details).
			WithParams(append([]string{resource, identifier}, params...)...),
	}
}
