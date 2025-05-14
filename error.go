package main

import (
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int
	Err        string
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprint(" ", e.StatusCode, " - ", e.Err, " - ", e.Message)
}

func NewBadRequestError(message string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Err:        "ERROR_BAD_REQUEST",
		Message:    message,
	}
}

func NewNotFoundError(message string) *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERROR_NOT_FOUND",
		Message:    message,
	}
}

func NewInternalServerError(err error) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Err:        "ERROR_INTERNAL_SERVER",
		Message:    err.Error(),
	}
}
