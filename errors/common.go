// Package errors defines common, reusable sentinel errors used across
// the Anedya Go SDK.
//
// These errors represent generic failure categories that are not tied
// to any specific API domain and are intended to be wrapped by higher-level
// errors for additional context.
package errors

import "errors"

// Generic SDK errors used by ALL APIs.
var (
	// ErrRequestEncodeFailed indicates a failure while encoding
	// or serializing the request payload.
	ErrRequestEncodeFailed = errors.New("request encode failed")

	// ErrRequestBuildFailed indicates a failure while constructing
	// the HTTP request.
	ErrRequestBuildFailed = errors.New("request build failed")

	// ErrRequestFailed indicates a failure during HTTP request execution.
	ErrRequestFailed = errors.New("request failed")

	// ErrResponseDecodeFailed indicates a failure while decoding
	// the API response body.
	ErrResponseDecodeFailed = errors.New("response decode failed")

	// ErrUnauthorized indicates an authentication or authorization failure.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrUnknown indicates an unclassified or unexpected error.
	ErrUnknown = errors.New("unknown error")
)
