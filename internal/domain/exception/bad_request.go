package exception

import (
	"net/http"
)

type BadRequestException struct {
	BaseError
}

func NewBadRequestException(msg string, code string, meta map[string]any) *BadRequestException {
	if meta == nil {
		meta = make(map[string]any)
	}

	return &BadRequestException{
		BaseError: BaseError{
			errorCode:  code,
			message:    msg,
			httpStatus: http.StatusBadRequest, // 400
			details:    meta,
		},
	}
}

func NewInvalidRequestBodyException(err error) *BadRequestException {
	return NewBadRequestException(
		"Invalid request body: "+err.Error(),
		"INVALID_JSON_FORMAT",
		nil,
	)
}

func NewValidationFailedException(details map[string]any) *BadRequestException {
	return NewBadRequestException(
		"Input validation failed. Please review the details for specific field issues.",
		"INPUT_VALIDATION_FAILED",
		details,
	)
}

func NewVerificationCodeExpired(err error) *BadRequestException {
	ex := NewBadRequestException(
		"Verification code expired or not found.",
		"VERIFY_CODE_EXPIRED",
		nil,
	)
	if err != nil {
		ex.Wrap(err)
	}
	return ex
}

func NewInvalidVerificationCode() *BadRequestException {
	return NewBadRequestException(
		"Invalid verification code.",
		"VERIFY_CODE_INVALID",
		nil,
	)
}

func (e *BadRequestException) ClientError() {}