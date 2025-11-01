package exception

import (
	"fmt"
	"net/http"
)

type InternalServerException struct {
	BaseError
}

func NewInternalServerException(clientMsg string, code string, err error) *InternalServerException {
	e := &InternalServerException{
		BaseError: BaseError{
			errorCode:  code,
			message:    clientMsg,
			httpStatus: http.StatusInternalServerError, // 500
			details:    nil,
		},
	}
	if err != nil {
		e.Wrap(err)
	}
	return e
}

func NewRepositoryError(operation string, err error) *InternalServerException {
	return NewInternalServerException(
		fmt.Sprintf("Failed to complete repository operation: %s", operation),
		"DB_OPERATION_FAILED",
		err,
	)
}

func NewHashedPasswordError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to hash password",
		"SYSTEM_HASH_FAIL",
		err,
	)
}

func NewJWTError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to generate JWT token",
		"SYSTEM_JWT_FAIL",
		err,
	)
}

func NewVerificationError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to verify email",
		"SYSTEM_VERIFY_FAIL",
		err,
	)
}

func NewEmailError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to send email",
		"SYSTEM_EMAIL_FAIL",
		err,
	)
}

func NewDBLoginError(err error) *InternalServerException {	
	return NewInternalServerException(
		"Database failure during login",
		"DB_LOGIN_FAIL",
		err,
	)
}

func NewRepositoryVerificationError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to verify email",
		"DB_VERIFY_FAIL",
		err,
	)
}

func NewRepositoryUpdateError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to update user",
		"DB_UPDATE_FAIL",
		err,
	)
}

func NewVerificationCodeGenerationError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to generate verification code",
		"VERIFY_CODE_GEN_FAIL",
		err,
	)
}

func NewContextCastError(err error) *InternalServerException {
	return NewInternalServerException(
		"Failed to cast context",
		"CONTEXT_CAST_FAIL",
		err,
	)
}
