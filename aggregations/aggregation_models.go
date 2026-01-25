// Package aggregations provides types and configurations
// for defining time-based data aggregations within the
// Anedya platform.
package aggregations

// AggregationType represents the type of aggregation
// operation to be computed on time-series data.
type AggregationType string

const (
	// AggregationSum computes the sum of all values
	// over the specified interval.
	AggregationSum AggregationType = "sum"

	// AggregationAvg computes the average value
	// over the specified interval.
	AggregationAvg AggregationType = "avg"

	// AggregationMedian computes the median value
	// over the specified interval.
	AggregationMedian AggregationType = "median"

	// AggregationMin computes the minimum value
	// over the specified interval.
	AggregationMin AggregationType = "min"

	// AggregationMax computes the maximum value
	// over the specified interval.
	AggregationMax AggregationType = "max"

	// AggregationDiff computes the difference between
	// the first and last value in the interval.
	AggregationDiff AggregationType = "diff"

	// AggregationDeltaSum computes the delta sum across
	// the interval while accounting for counter resets.
	AggregationDeltaSum AggregationType = "deltasum"

	// AggregationStdDev computes the standard deviation
	// of values over the specified interval.
	AggregationStdDev AggregationType = "stddev"
)

// IntervalMeasure represents the unit of time used
// for defining aggregation intervals.
type IntervalMeasure string

const (
	// MeasureYear represents yearly aggregation intervals.
	MeasureYear IntervalMeasure = "year"

	// MeasureMonth represents monthly aggregation intervals.
	MeasureMonth IntervalMeasure = "month"

	// MeasureWeek represents weekly aggregation intervals.
	MeasureWeek IntervalMeasure = "week"

	// MeasureDay represents daily aggregation intervals.
	MeasureDay IntervalMeasure = "day"

	// MeasureHour represents hourly aggregation intervals.
	MeasureHour IntervalMeasure = "hour"

	// MeasureMinute represents minute-based aggregation intervals.
	MeasureMinute IntervalMeasure = "minute"
)

// FilterType represents the type of node filtering
// applied during aggregation.
type FilterType string

const (
	// FilterInclude includes only the specified nodes
	// in the aggregation computation.
	FilterInclude FilterType = "include"

	// FilterExclude excludes the specified nodes
	// from the aggregation computation.
	FilterExclude FilterType = "exclude"
)

// AggregationConfig holds the complete configuration
// required to perform a data aggregation.
type AggregationConfig struct {
	// Aggregation defines the aggregation method and options.
	Aggregation AggregationObject `json:"aggregation"`

	// Interval defines the time-based aggregation window.
	Interval IntervalObject `json:"interval"`

	// ResponseOptions defines how the aggregated response
	// should be formatted.
	ResponseOptions ResponseOptions `json:"responseOptions,omitempty"`

	// Filter defines optional node filtering rules
	// applied before aggregation.
	Filter *FilterObject `json:"filter,omitempty"`
}

// AggregationObject defines the aggregation operation
// to be performed on the data.
type AggregationObject struct {
	// Compute specifies the type of aggregation
	// (sum, avg, min, max, etc.).
	Compute AggregationType `json:"compute"`

	// ForEachNode determines whether aggregation is computed
	// individually per node or combined across all nodes.
	ForEachNode bool `json:"forEachNode,omitempty"`
}

// IntervalObject defines the size and unit of the
// aggregation time window.
type IntervalObject struct {
	// Measure specifies the unit of time
	// (minute, hour, day, etc.).
	Measure IntervalMeasure `json:"measure"`

	// Interval specifies the number of units
	// for the aggregation window (e.g., 5 for "5 minutes").
	Interval int `json:"interval"`
}

// ResponseOptions defines options for formatting
// the aggregation response.
type ResponseOptions struct {
	// Timezone specifies the timezone used for
	// timestamps in the response (default: UTC).
	Timezone string `json:"timezone,omitempty"`
}

// FilterObject defines node-level filtering rules
// applied before aggregation.
type FilterObject struct {
	// Nodes is the list of node IDs used for filtering.
	Nodes []string `json:"nodes,omitempty"`

	// Type specifies whether nodes are included or excluded.
	// Allowed values are "include" or "exclude".
	Type FilterType `json:"type,omitempty"`
}

// AggregateDataPoint represents a single aggregated
// value for a specific time interval.
type AggregateDataPoint struct {
	// Timestamp indicates the start of the aggregation
	// interval (Unix milliseconds).
	Timestamp int64 `json:"timestamp"`

	// Aggregate holds the computed aggregation value.
	Aggregate float64 `json:"aggregate"`
}
