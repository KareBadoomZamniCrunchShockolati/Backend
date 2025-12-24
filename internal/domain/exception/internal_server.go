package exception

import "net/http"

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

type InternalServerException struct {
	*BaseError
}

func NewInternalServerException(code string, details map[string]any, originalErr error, params ...string) *InternalServerException {
	if details == nil {
		details = map[string]any{} 
	}

	be := NewBaseError(code, http.StatusInternalServerError, details).WithParams(params...)
	if originalErr != nil {
		be = be.Wrap(originalErr)
	}

	return &InternalServerException{BaseError: be}
}

func NewInternalServerExceptionWithErr(code string, originalErr error) *InternalServerException {
	return NewInternalServerException(code, nil, originalErr)
}

func NewRepositoryError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeDBOperationFailed, nil, err)
}

func NewHashedPasswordError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeHashFail, nil, err)
}

func NewJWTError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeJWTFail, nil, err)
}

func NewVerificationError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeVerifyFail, nil, err)
}

func NewEmailError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeEmailFail, nil, err)
}

func NewDBLoginError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeDBLoginFail, nil, err)
}

func NewRepositoryVerificationError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeDBVerifyFail, nil, err)
}

func NewRepositoryUpdateError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeDBUpdateFail, nil, err)
}

func NewVerificationCodeGenerationError(err error) *InternalServerException {
	return NewInternalServerException(ErrorTypeVerifyCodeGenFail, nil, err)
}

func NewContextCastError(err error) *InternalServerException {
	return NewInternalServerException(
		ErrorTypeContextCastFail,
		map[string]any{
			"reason": "context_cast_failed",
		},
		err,
	)
}
