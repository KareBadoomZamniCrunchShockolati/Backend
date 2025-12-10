package exception

import "fmt"

type Error interface {
	error
	Code() string
	HTTPStatus() int
	Details() map[string]any
	Unwrap() error
}

type ClientError interface {
	Error
	ClientError()
}

type BaseError struct {
	errorCode  string
	message    string
	httpStatus int
	details    map[string]any
	internal   error
}

func (e *BaseError) Error() string {
	if e.internal != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.errorCode, e.message, e.internal)
	}
	return fmt.Sprintf("%s: %s", e.errorCode, e.message)
}

func (e *BaseError) Code() string            { return e.errorCode }
func (e *BaseError) HTTPStatus() int         { return e.httpStatus }
func (e *BaseError) Details() map[string]any { return e.details }
func (e *BaseError) Unwrap() error           { return e.internal }
func (e *BaseError) Wrap(err error) *BaseError {
	e.internal = err
	return e
}
func NewBaseError(code, message string, httpStatus int, details map[string]any) *BaseError {
	return &BaseError{
		errorCode:  code,
		message:    message,
		httpStatus: httpStatus,
		details:    details,
	}
}
