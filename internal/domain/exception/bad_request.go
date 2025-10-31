package exception

import (
	"net/http"
)

type BadRequestException struct {
	BaseError
}

func NewBadRequestException(msg string, code string, meta map[string]any) BadRequestException {
	if meta == nil {
		meta = make(map[string]any)
	}
	
	return BadRequestException{
		BaseError: BaseError{
			errorCode:  code,
			message:    msg,
			httpStatus: http.StatusBadRequest, // 400
			details:    meta,
		},
	}
}

func (e BadRequestException) ClientError() {}