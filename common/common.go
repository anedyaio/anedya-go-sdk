// Package common contains shared models and utilities
// used across multiple API modules in the SDK.
package common

// BaseResponse represents the common structure
// returned by all API responses.
//
// It is embedded or extended by module-specific
// response types to ensure consistent error handling
// across the SDK.
type BaseResponse struct {
	// Success indicates whether the API request
	// was processed successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code
	// returned by the API for programmatic handling.
	ReasonCode string `json:"reasonCode,omitempty"`
}
