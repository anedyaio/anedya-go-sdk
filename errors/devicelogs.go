package errors

import "errors"

// Device Logs API – API level errors
var (
	// ErrLogsRetentionLimit indicates that the requested logs
	// exceed the maximum retention period supported by the platform.
	//
	// Currently, device logs are preserved for a maximum of 6 months.
	ErrLogsRetentionLimit = errors.New(
		"logs are preserved for a maximum of 6 months",
	)
)
