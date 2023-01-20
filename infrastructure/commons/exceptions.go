package commons

import "errors"

var (
	// ErrInvalidCredentials is thrown when the user credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInternalServerError is thrown when the server encounters an error
	ErrInternalServerError = errors.New("internal server error")
	// ErrUserAlreadyExists is thrown when the user already exists
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrEmptyInput is thrown when the input is empty
	ErrEmptyInput = errors.New("empty input")
	// ErrValidationFailed is thrown when the input validation is failed
	ErrValidationFailed = errors.New("validation failed")
)