package utils

import (
	"errors"
	"fmt"
)

// Errors that can be used to check with error.Is()
var (
	// Validation errors
	ErrInvalidInput    = errors.New("invalid input provided")
	ErrInvalidType     = errors.New("invalid type provided")
	ErrMissingRequired = errors.New("required data is missing")

	// API errors
	ErrResourceNotFound = errors.New("resource not found")
	ErrAccessDenied     = errors.New("access denied")
	ErrUnauthorized     = errors.New("unauthorized request")
	ErrMalformedRequest = errors.New("malformed request")
	ErrUnkown           = errors.New("unkown error occured")

	// Netowrk errors
	ErrRequestFailed = errors.New("request failed")

	// VariableManagement-specific sentinels
	ErrVariableNotFound = errors.New("variable not found")

	// NodeManagement- specific sentinels
)

// AnedyaError struct
// This holds the extra details from the API and it wraps a sentinel error
type AnedyaError struct {
	Code    string // The API string code
	Message string // The detailed API message
	Err     error  // The sentinel error
}

func (e *AnedyaError) Error() string {
	// If we have a message, use it otherwise use sentinel's message
	if e.Message != "" {
		return fmt.Sprintf("anedya api error: %s (code: %s)", e.Message, e.Code)
	}
	if e.Err != nil {
		return fmt.Sprintf("anedya api error: %s (code: %s)", e.Err.Error(), e.Code)
	}
	return fmt.Sprintf("anedya api error: unkown error (code: %s)", e.Code)
}

// Unwrap implementation
// This allows error.Is() to work with AnedyaError
func (e *AnedyaError) Unwrap() error {
	return e.Err
}

// Code Map (API Code -> Sentinel errors)
var CodeMap = map[string]error{
	// Generic SDK errors
	"sdk::invalidinput":    ErrInvalidInput,
	"sdk::missingrequired": ErrMissingRequired,
	"sdk::invalidtype":     ErrInvalidType,

	// Generic API errors
	"auth::accessdenied": ErrAccessDenied,

	// Variable errors
	"variable::notfound": ErrVariableNotFound,

	// Token errors
}

// GetErrorCode extracts the API error code from an error
func GetErrorCode(err error) string {
	if anedyaErr, ok := err.(*AnedyaError); ok {
		return anedyaErr.Code
	}
	return ""
}


