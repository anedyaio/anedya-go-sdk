package valuestore

import (
	"fmt"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// AsString attempts to retrieve the stored value as a string.
//
// It performs a type assertion on the underlying interface{}.
// If the value is not a string, it returns an empty string and an ErrTypeMismatch error.
func (r *GetValueResponse) AsString() (string, error) {
	v, ok := r.Value.(string)
	if ok {
		return v, nil
	}
	return "", &errors.AnedyaError{
		Message: fmt.Sprintf("expected string, got %T", r.Value),
		Err:     errors.ErrTypeMismatch,
	}
}

// AsFloat attempts to retrieve the stored value as a float64.
//
// If the value is not a number, it returns 0 and an ErrTypeMismatch error.
func (r *GetValueResponse) AsFloat() (float64, error) {
	v, ok := r.Value.(float64)
	if ok {
		return v, nil
	}
	return 0, &errors.AnedyaError{
		Message: fmt.Sprintf("value is not a number, got type: %T", r.Value),
		Err:     errors.ErrTypeMismatch,
	}
}

// AsBool attempts to retrieve the stored value as a boolean.
//
// It performs a type assertion on the underlying interface{}.
// If the value is not a boolean, it returns false and an ErrTypeMismatch error.
func (r *GetValueResponse) AsBool() (bool, error) {
	v, ok := r.Value.(bool)
	if ok {
		return v, nil
	}
	return false, &errors.AnedyaError{
		Message: fmt.Sprintf("value is not a boolean, got type: %T", r.Value),
		Err:     errors.ErrTypeMismatch,
	}
}
