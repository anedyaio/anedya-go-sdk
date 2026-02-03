// Package errors defines reusable sentinel errors used across
// the Anedya Go SDK. These errors represent validation failures,
// missing required fields, and API-level command issues.
package errors

import "errors"

var (
	// ErrInvalidCommandID indicates that the provided command ID
	// is empty, malformed, or does not meet expected format rules.
	ErrInvalidCommandID = errors.New("invalid command id")

	// ErrInvalidLimit indicates that the pagination limit
	// is outside the allowed range.
	ErrInvalidLimit = errors.New("invalid limit")

	// ErrInvalidOffset indicates that the pagination offset
	// is negative or invalid.
	ErrInvalidOffset = errors.New("invalid offset")

	// ErrInvalidCommandType indicates that the command data type
	// is neither supported nor recognized.
	ErrInvalidCommandType = errors.New("invalid command type")

	// ErrDataRequired indicates that the command payload data
	// was not provided but is mandatory.
	ErrDataRequired = errors.New("data is required")

	// ErrCommandRequired indicates that the command name
	// or identifier was not provided.
	ErrCommandRequired = errors.New("command is required")

	// ErrInvalidExpiry indicates that the command expiry
	// value is negative or otherwise invalid.
	ErrInvalidExpiry = errors.New("invalid command expiry")

	// ErrUnknownDataType indicates that the API returned
	// an unsupported or unknown data type.
	ErrUnknownDataType = errors.New("unknown command data type")

	// ErrCommandNotFound indicates that the command
	// does not exist or is no longer pending.
	ErrCommandNotFound = errors.New("command not found or not pending")

	// ErrInvalidFilter indicates that the provided
	// command filter configuration is invalid.
	ErrInvalidFilter = errors.New("invalid command filter")
)
