package exception

import "fmt"

type Error interface {
	error
	Code() string
	HTTPStatus() int
	Details() map[string]any
	Params() []string
	Unwrap() error
}

type ClientError interface {
	Error
	ClientError()
}

type BaseError struct {
	errorCode  string
	httpStatus int
	details    map[string]any
	params     []string
	internal   error
}

func (e *BaseError) Error() string {
	if e.internal != nil {
		return fmt.Sprintf("%s (status=%d) details=%v cause=%v", e.errorCode, e.httpStatus, e.details, e.internal)
	}
	return fmt.Sprintf("%s (status=%d) details=%v", e.errorCode, e.httpStatus, e.details)
}

func (e *BaseError) Code() string            { return e.errorCode }
func (e *BaseError) HTTPStatus() int         { return e.httpStatus }
func (e *BaseError) Details() map[string]any { return e.details }
func (e *BaseError) Params() []string        { return e.params }
func (e *BaseError) Unwrap() error           { return e.internal }

func (e *BaseError) Wrap(err error) *BaseError {
	e.internal = err
	return e
}

func (e *BaseError) WithParams(params ...string) *BaseError {
	e.params = params
	return e
}

func NewBaseError(code string, httpStatus int, details map[string]any) *BaseError {
	return &BaseError{
		errorCode:  code,
		httpStatus: httpStatus,
		details:    details,
	}
}
