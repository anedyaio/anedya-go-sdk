// Package errors defines SDK and API error types used by the
// Anedya Go SDK.
package errors

import (
	"fmt"
)

// AnedyaError represents a structured SDK or API error.
//
// It wraps a sentinel error for programmatic checks and includes
// a human-readable message for debugging or logging.
type AnedyaError struct {
	// Message is a human-readable error description.
	Message string

	// Err is the underlying sentinel error.
	Err error
}

// Error implements the error interface.
func (e *AnedyaError) Error() string {
	return fmt.Sprintf("anedya api error: %s: %v", e.Message, e.Err)
}

// Unwrap allows errors.Is to work with AnedyaError.
func (e *AnedyaError) Unwrap() error {
	return e.Err
}

// codeMap maps API reason codes to SDK sentinel errors.
var codeMap = map[string]error{
	// node errors
	"node::devidexists":          ErrNodeDeviceIDExists,
	"node::childexists":          ErrNodeChildExists,
	"node::uniquealiasviolation": ErrNodeUniqueAliasViolation,
	"node::uniquechildviolation": ErrNodeUniqueChildViolation,
	"node::childnotfound":        ErrNodeChildNotFound,
	"node::invalidparentid":      ErrNodeInvalidParentID,
	"node::invalidchildid":       ErrNodeInvalidChildID,
	"node::nodenotfound":         ErrNodeNotFound,
	"node::devicenotfound":       ErrNodeDeviceNotFound,
	"node::invaliduuid":          ErrNodeInvalidUUID,

	// Data API errors
	"data::variablenotfound": ErrVariableNotFound,
	"data::invalidnodeid":    ErrInvalidNodeID,

	"variable::namerequired":     ErrVariableNameRequired,
	"variable::variablerequired": ErrVariableRequired,
	"variable::typerequired":     ErrVariableTypeRequired,

	//auth
	"auth::accessdenied":  ErrAccessDenied,
	"auth::tokennotfound": ErrTokenNotFound,
}

// GetError converts an API reason code and message into an AnedyaError.
func GetError(code, message string) error {
	sentinel, ok := codeMap[code]
	if !ok {
		sentinel = ErrUnknown
	}

	return &AnedyaError{
		Message: message,
		Err:     sentinel,
	}
}
