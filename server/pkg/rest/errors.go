package rest

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status     int    `json:"status,omitempty"`
	StatusText string `json:"statusText,omitempty"`
	Message    string `json:"error,omitempty"`
}

func NewError(status int, message string) *Error {
	return &Error{
		Status:     status,
		StatusText: http.StatusText(status),
		Message:    message,
	}
}

func NewBadRequestError(message string) *Error {
	return NewError(http.StatusBadRequest, message)
}

func NewNotFoundError(message string) *Error {
	return NewError(http.StatusNotFound, message)
}

func NewConflictError(message string) *Error {
	return NewError(http.StatusConflict, message)
}

func NewUnauthorizedError(message string) *Error {
	return NewError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *Error {
	return NewError(http.StatusForbidden, message)
}

func NewInternalServerError() *Error {
	return NewError(http.StatusInternalServerError, "Internal Server Error Occurred")
}

func (e Error) Error() string {
	return fmt.Sprintf("[%v]:%s - %s", e.Status, e.StatusText, e.Message)
}

func (e Error) Response() (int, interface{}) {
	return e.Status, e
}
