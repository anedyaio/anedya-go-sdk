package errors

import "errors"

var (
	// ErrInvalidAggregationMethod is returned when an unsupported
	// aggregation method is provided in a request.
	ErrInvalidAggregationMethod = errors.New("invalid aggregation method")

	// ErrAggregationMethodRequired is returned when the aggregation
	// method is missing but is required for the operation.
	ErrAggregationMethodRequired = errors.New("aggregation method is required")

	// ErrInvalidIntervalMeasure is returned when the interval
	// measurement unit (e.g., seconds, minutes, hours) is invalid.
	ErrInvalidIntervalMeasure = errors.New("invalid interval measure")

	// ErrInvalidInterval is returned when the provided interval
	// value is invalid or out of the allowed range.
	ErrInvalidInterval = errors.New("invalid interval")

	// ErrInvalidTimezone is returned when an unsupported or
	// malformed timezone is provided.
	ErrInvalidTimezone = errors.New("invalid timezone")

	// ErrInvalidFilterType is returned when an unsupported
	// filter type is used in a request.
	ErrInvalidFilterType = errors.New("invalid filter type")
)
