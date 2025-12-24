package exception

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	ErrorTypeInvalidJSONFormat       = "INVALID_REQUEST_BODY"
	ErrorTypeInputValidationFailed   = "INPUT_VALIDATION_FAILED"
	ErrorTypeVerificationCodeExpired = "VERIFY_CODE_EXPIRED"
	ErrorTypeVerificationCodeInvalid = "VERIFY_CODE_INVALID"
)

type BadRequestException struct {
	*BaseError
}

func (e *BadRequestException) ClientError() {}

func NewBadRequestException(code string, meta map[string]any, params ...string) *BadRequestException {
	if meta == nil {
		meta = map[string]any{}
	}
	return &BadRequestException{
		BaseError: NewBaseError(code, http.StatusBadRequest, meta).WithParams(params...),
	}
}

func NewInvalidRequestBodyException(err error) *BadRequestException {
	ex := NewBadRequestException(ErrorTypeInvalidJSONFormat, nil)
	if err != nil {
		ex.Wrap(err)
	}
	return ex
}

func NewValidationFailedException(details map[string]any) *BadRequestException {
	return NewBadRequestException(ErrorTypeInputValidationFailed, details)
}

func NewValidationException(err error) Error {
	if err == nil {
		return nil
	}
	var jsonErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) {
		return NewInvalidRequestBodyException(err)
	}
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		fields := make([]map[string]any, 0, len(verr))
		for _, fe := range verr {
			fields = append(fields, map[string]any{
				"field": fe.Field(), 
				"tag":   fe.Tag(),
				"param": fe.Param(),
			})
		}
		return NewValidationFailedException(map[string]any{"fields": fields})
	}
	return NewInvalidRequestBodyException(err)
}

func NewVerificationCodeExpired(err error) *BadRequestException {
	ex := NewBadRequestException(ErrorTypeVerificationCodeExpired, nil)
	if err != nil {
		ex.Wrap(err)
	}
	return ex
}

func NewInvalidVerificationCode() *BadRequestException {
	return NewBadRequestException(ErrorTypeVerificationCodeInvalid, nil)
}
