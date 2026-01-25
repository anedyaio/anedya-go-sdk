// Package errors defines reusable sentinel errors used across
// the Anedya Go SDK for validation and configuration failures.
package errors

import "errors"

var (
	// ErrInvalidAggregationMethod is returned when an unsupported
	// aggregation method is provided.
	ErrInvalidAggregationMethod = errors.New("invalid aggregation method")

	// ErrAggregationMethodRequired is returned when the aggregation
	// method field is missing or empty.
	ErrAggregationMethodRequired = errors.New("aggregation method is required")

	// ErrInvalidIntervalMeasure is returned when an unsupported
	// interval measure (year, month, day, etc.) is provided.
	ErrInvalidIntervalMeasure = errors.New("invalid interval measure")

	// ErrInvalidInterval is returned when the interval value
	// is zero or negative.
	ErrInvalidInterval = errors.New("invalid interval")

	// ErrInvalidTimezone is returned when an invalid or
	// unsupported timezone is provided.
	ErrInvalidTimezone = errors.New("invalid timezone")

	// ErrInvalidFilterType is returned when an unsupported
	// filter type is specified.
	ErrInvalidFilterType = errors.New("invalid filter type")

	// ErrFilterNodesRequired is returned when filter nodes
	// are required but not provided.
	ErrFilterNodesRequired = errors.New("filter nodes required")

	// ErrAggregationComputeRequired is returned when the
	// aggregation compute type is missing or empty.
	ErrAggregationComputeRequired = errors.New("aggregation compute type required")
)
