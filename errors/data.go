package errors

import "errors"

// Generic validation errors
var (
	ErrRequestNil       = errors.New("request is nil")
	ErrNodesEmpty       = errors.New("nodes list is empty")
	ErrInvalidNode      = errors.New("invalid node")
	ErrInvalidTimeRange = errors.New("invalid from/to time range")
	ErrInvalidOrder     = errors.New("invalid order value")
	ErrInvalidTimestamp = errors.New("invalid timestamp")
)

// Data API – API level errors
var (
	ErrVariableNotFound = errors.New("variable not found")
	ErrInvalidNodeID    = errors.New("invalid node id")
)
