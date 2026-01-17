// Package errors defines variable-related validation and API errors
// used by the Variable Management APIs in the Anedya Go SDK.
package errors

import "errors"

// Variable Management validation errors.
//
// These sentinel errors are returned when a Variable API request
// fails basic validation before being sent to the Anedya API.
var (
	// ErrVariableNameRequired is returned when the variable name
	// field is missing or empty in the request.
	ErrVariableNameRequired = errors.New("variable name is required")

	// ErrVariableRequired is returned when the variable identifier
	// (key/path) is missing or empty in the request.
	ErrVariableRequired = errors.New("variable is required")

	// ErrVariableTypeRequired is returned when the variable type
	// is missing or not one of the supported values.
	//
	// Supported types:
	//   - geo
	//   - float
	ErrVariableTypeRequired = errors.New("variable type is required")
)
