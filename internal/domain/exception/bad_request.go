package exception

import (
	"net/http"
)

const (
	ErrorTypeInvalidJSONFormat       = "INVALID_JSON_FORMAT"
	ErrorTypeInputValidationFailed   = "INPUT_VALIDATION_FAILED"
	ErrorTypeVerificationCodeExpired = "VERIFY_CODE_EXPIRED"
	ErrorTypeVerificationCodeInvalid = "VERIFY_CODE_INVALID"
)

type BadRequestException struct {
	*BaseError
}

func NewBadRequestException(msg string, code string, meta map[string]any) *BadRequestException {
	if meta == nil {
		meta = make(map[string]any)
	}
	return &BadRequestException{
		BaseError: NewBaseError(
			code,
			msg,
			http.StatusBadRequest, // 400
			meta,
		),
	}
}

func NewInvalidRequestBodyException(err error) *BadRequestException {
	return NewBadRequestException(
		"Invalid request body: "+err.Error(),
		ErrorTypeInvalidJSONFormat,
		nil,
	)
}

func NewValidationFailedException(details map[string]any) *BadRequestException {
	return NewBadRequestException(
		"Input validation failed. Please review the details for specific field issues.",
		ErrorTypeInputValidationFailed,
		details,
	)
}

func NewVerificationCodeExpired(err error) *BadRequestException {
	ex := NewBadRequestException(
		"Verification code expired or not found.",
		ErrorTypeVerificationCodeExpired,
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
		ErrorTypeVerificationCodeInvalid,
		nil,
	)
}

func (e *BadRequestException) ClientError() {}
