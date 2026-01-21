// Package errors defines validation and API errors
// used by the Access Token APIs in the Anedya Go SDK.
package errors

import "errors"

var (
	// ErrExpiryRequired indicates that the token expiry value (ttlSec) was not provided or is invalid.
	ErrExpiryRequried = errors.New("expiry value is required")

	// ErrPolicyRequired indicates that the access policy object is missing.
	ErrPolicyRequired = errors.New("policy is missing")

	// ErrInvalidPermission indicates that one or more permissions
	// provided in the policy.allow field are invalid or unsupported.
	ErrInavalidPermission = errors.New("permissions are invalid")

	// ErrTokenIdRequired indicates that the tokenId field
	// is missing or empty when it is required by the API.
	ErrTokenIdRequired = errors.New("TokenID value is required")
)
