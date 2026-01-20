package errors

import "errors"

var (
	ErrExpiryRequried     = errors.New("expiry value is required")
	ErrPolicyRequired     = errors.New("policy is missing")
	ErrInavalidPermission = errors.New("permissions are invalid")
	ErrTokenIdRequired    = errors.New("TokenID value is required")
)
