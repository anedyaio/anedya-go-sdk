package errors

import "errors"

var (
	// ErrHealthLimitExceeded is returned when the provided
	// lastContactThreshold exceeds the maximum allowed
	// duration for health status checks.
	ErrHealthLimitExceeded = errors.New("health lastContactThreshold exceeded")
)
