package exception

import (
	"net/http"
)

const (
	ErrorTypeDBOperationFailed = "DB_OPERATION_FAILED"
	ErrorTypeHashFail          = "SYSTEM_HASH_FAIL"
	ErrorTypeJWTFail           = "SYSTEM_JWT_FAIL"
	ErrorTypeVerifyFail        = "SYSTEM_VERIFY_FAIL"
	ErrorTypeEmailFail         = "SYSTEM_EMAIL_FAIL"
	ErrorTypeDBLoginFail       = "DB_LOGIN_FAIL"
	ErrorTypeDBVerifyFail      = "DB_VERIFY_FAIL"
	ErrorTypeDBUpdateFail      = "DB_UPDATE_FAIL"
	ErrorTypeVerifyCodeGenFail = "VERIFY_CODE_GEN_FAIL"
	ErrorTypeContextCastFail   = "CONTEXT_CAST_FAIL"
)

// InternalServerException represents a 500 Internal Server Error
type InternalServerException struct {
	*BaseError
}

// NewInternalServerException creates a new InternalServerException
func NewInternalServerException(code string, message string, originalErr error) *InternalServerException {
	return &InternalServerException{
		BaseError: NewBaseError(
			code,
			message, // descriptive message
			http.StatusInternalServerError,
			nil,
		).Wrap(originalErr),
	}
}

// Convenience constructors
func NewRepositoryError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeDBOperationFailed,
		"Failed to complete repository operation",
		err,
	)
}

func NewHashedPasswordError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeHashFail,
		"Failed to hash password",
		err,
	)
}

func NewJWTError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeJWTFail,
		"Failed to generate JWT token",
		err,
	)
}

func NewVerificationError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeVerifyFail,
		"Failed to verify email",
		err,
	)
}

func NewEmailError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeEmailFail,
		"Failed to send email",
		err,
	)
}

func NewDBLoginError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeDBLoginFail,
		"Database failure during login",
		err,
	)
}

func NewRepositoryVerificationError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeDBVerifyFail,
		"Failed to verify email in repository",
		err,
	)
}

func NewRepositoryUpdateError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeDBUpdateFail,
		"Failed to update user in repository",
		err,
	)
}

func NewVerificationCodeGenerationError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeVerifyCodeGenFail,
		"Failed to generate verification code",
		err,
	)
}

func NewContextCastError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeContextCastFail,
		"Failed to cast context",
		err,
	)
}
