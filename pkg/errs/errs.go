package errs

import (
	"fmt"
	"net/http"
)

type ClientError interface {
	error
	Status() int
	Message() string
}

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.Resource)
}

func (e *NotFoundError) Status() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Message() string {
	return fmt.Sprintf("The resource %s was not found.", e.Resource)
}

type ConflictError struct {
	MessageValue string 
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict detected: %s", e.MessageValue)
}

func (e *ConflictError) Status() int {
	return http.StatusConflict 
}

func (e *ConflictError) Message() string {
	return e.Error()
}

type InternalServerError struct {
	Err error 
}

func (e *InternalServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal server error: %v", e.Err)
	}
	return "internal server error"
}

func (e *InternalServerError) Status() int {
	return http.StatusInternalServerError 
}

func (e *InternalServerError) Message() string {
	return "An unexpected server error occurred." 
}

type UnAuthorizedError struct {
	MessageValue string
}

func (e *UnAuthorizedError) Error() string {
	return fmt.Sprintf("unauthorized access: %s", e.MessageValue)
}

func (e *UnAuthorizedError) Status() int {
	return http.StatusUnauthorized
}

func (e *UnAuthorizedError) Message() string {
	return fmt.Sprintf("You are not authorized to access this resource: %s", e.MessageValue)
}

type BadRequestError struct {
	MessageValue string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.MessageValue)
}

func (e *BadRequestError) Status() int {
	return http.StatusBadRequest
}

func (e *BadRequestError) Message() string {
	return fmt.Sprintf("Invalid request: %s", e.MessageValue)
}
