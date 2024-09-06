package errors

import (
	"net/http"
)

// Error represents a custom error with an HTTP status code.
type Error struct {
	Code    int
	Message string
}

// New creates a new Error instance.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// HandleError writes the error response to the HTTP response writer.
func HandleError(w http.ResponseWriter, err *Error) {
	http.Error(w, err.Message, err.Code)
}

// Example error instances
var (
	ErrBadRequest          = New(http.StatusBadRequest, "Bad Request")
	ErrForbidden           = New(http.StatusForbidden, "Forbidden")
	ErrNotFound            = New(http.StatusNotFound, "Not Found")
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal Server Error")
)
