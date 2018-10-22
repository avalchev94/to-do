package api

import "errors"

// Error is a struct for easy working with API errors
type Error struct {
	Err string `json:"error"`
}

var (
	// Unauthorized - the user is not logged.
	Unauthorized = Error{"No authorized user"}
	// AlreardyAuth - a logged user is trying to login again.
	AlreardyAuth = Error{"Already authorized"}
	// IncorrectParameter - a generic error for routes with url parameter(user/:id)
	IncorrectParameter = Error{"Incorrect parameter"}
	// ResourceNotFound - when the request is correct, but there is no such data in DB.
	ResourceNotFound = Error{"Resource is not found"}
)

// NewError is a simple function for creating an api.Error
// from the builtin error type.
func NewError(err error) Error {
	return Error{err.Error()}
}

// Error converts api.Error to builtin type error.
func (e Error) Error() error {
	return errors.New(e.Err)
}
