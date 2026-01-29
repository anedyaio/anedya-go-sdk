// Package errors defines reusable sentinel errors used across
// the Anedya Go SDK for command-related validation and API failures.
package errors

import "errors"

var (
	// ErrInvalidCommandID is returned when an invalid or
	// malformed command ID is provided.
	ErrInvalidCommandID = errors.New("invalid command id")

	// ErrInvalidLimit is returned when the provided limit
	// value is zero or negative.
	ErrInvalidLimit = errors.New("invalid limit")

	// ErrInvalidOffset is returned when the provided offset
	// value is negative or otherwise invalid.
	ErrInvalidOffset = errors.New("invalid offset")

	// ErrInvalidCommandType is returned when an unsupported
	// command type is specified.
	ErrInvalidCommandType = errors.New("invalid command type")

	// ErrDataRequired is returned when command data
	// is required but not provided.
	ErrDataRequired = errors.New("data is required")

	// ErrCommandRequired is returned when the command
	// name or identifier is missing.
	ErrCommandRequired = errors.New("command is required")

	// ErrInvalidExpiry is returned when the command
	// expiry timestamp is invalid.
	ErrInvalidExpiry = errors.New("invalid command expiry")

	// ErrUnknownDataType is returned when an unsupported
	// command data type is provided.
	ErrUnknownDataType = errors.New("unknown command data type")

	// ErrCommandNotFound is returned when the specified
	// command does not exist or is no longer pending.
	ErrCommandNotFound = errors.New("command not found or not pending")

	// ErrInvalidFilter is returned when the provided
	// command filter is invalid or malformed.
	ErrInvalidFilter = errors.New("invalid command filter")
)
