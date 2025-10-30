package exception

import "fmt"

type Error interface {
	error
	Code() string
	HTTPStatus() int
	Details() map[string]any
	Wrap(err error) Error 
	Unwrap() error      
}

type BaseError struct {
	errorCode  string
	message    string
	httpStatus int
	details    map[string]any
	internal   error 
}

func (e BaseError) Error() string { 
    if e.internal != nil {
        return fmt.Sprintf("%s (Internal: %v)", e.message, e.internal)
    }
    return e.message 
}

func (e BaseError) Code() string { return e.errorCode }
func (e BaseError) HTTPStatus() int { return e.httpStatus }
func (e BaseError) Details() map[string]any { return e.details }
func (e BaseError) Unwrap() error { return e.internal }
func (e BaseError) Wrap(err error) Error {
	e.internal = err
	return e
}