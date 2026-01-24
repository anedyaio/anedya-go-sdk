package valuestore

import (
	"encoding/json"
	"fmt"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// AsString attempts to retrieve the stored value as a string.
//
// It unmarshals the raw JSON bytes into a string.
// If the underlying JSON is not a string, it returns an ErrTypeMismatch error.
func (v *Value) AsString() (string, error) {
	var s string
	if err := json.Unmarshal(v.Value, &s); err != nil {
		return "", &errors.AnedyaError{
			Message: fmt.Sprintf("value is not a valid string (raw: %s)", string(v.Value)),
			Err:     errors.ErrTypeMismatch,
		}
	}
	return s, nil
}

// AsFloat attempts to retrieve the stored value as a float64.
//
// It unmarshals the raw JSON bytes into a float64.
// If the underlying JSON is not a number, it returns an ErrTypeMismatch error.
func (v *Value) AsFloat() (float64, error) {
	var f float64
	if err := json.Unmarshal(v.Value, &f); err != nil {
		return 0, &errors.AnedyaError{
			Message: fmt.Sprintf("value is not a valid number (raw: %s)", string(v.Value)),
			Err:     errors.ErrTypeMismatch,
		}
	}
	return f, nil
}

// AsBool attempts to retrieve the stored value as a boolean.
//
// It unmarshals the raw JSON bytes into a boolean.
// If the underlying JSON is not a boolean, it returns an ErrTypeMismatch error.
func (v *Value) AsBool() (bool, error) {
	var b bool
	if err := json.Unmarshal(v.Value, &b); err != nil {
		return false, &errors.AnedyaError{
			Message: fmt.Sprintf("value is not a valid boolean (raw: %s)", string(v.Value)),
			Err:     errors.ErrTypeMismatch,
		}
	}
	return b, nil
}
