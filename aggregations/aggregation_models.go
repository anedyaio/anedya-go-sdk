// Package aggregations provides configuration types and helpers
// for defining time-based and node-based data aggregations
// within the Anedya platform.
package aggregations

// AggregationType represents the type of aggregation
// computation to perform on data points.
type AggregationType string

const (
	// AggregationSum computes the sum of all values
	// within the given time interval.
	AggregationSum AggregationType = "sum"

	// AggregationAvg computes the average value
	// within the given time interval.
	AggregationAvg AggregationType = "avg"

	// AggregationMedian computes the median value
	// within the given time interval.
	AggregationMedian AggregationType = "median"

	// AggregationMin computes the minimum value
	// within the given time interval.
	AggregationMin AggregationType = "min"

	// AggregationMax computes the maximum value
	// within the given time interval.
	AggregationMax AggregationType = "max"

	// AggregationDiff computes the difference between
	// the first and last value in the interval.
	AggregationDiff AggregationType = "diff"

	// AggregationDeltaSum computes the difference between
	// the first and last value while accounting for
	// counter resets (useful for monotonically increasing counters).
	AggregationDeltaSum AggregationType = "deltasum"

	// AggregationStdDev computes the standard deviation
	// of values within the given interval.
	AggregationStdDev AggregationType = "stddev"
)

// IntervalMeasure represents the unit of time used
// to define aggregation intervals.
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

	// MeasureMinute represents minute-level aggregation intervals.
	MeasureMinute IntervalMeasure = "minute"
)

// FilterType represents the strategy used to filter nodes
// during aggregation.
type FilterType string

const (
	// FilterInclude includes only the specified nodes
	// in the aggregation.
	FilterInclude FilterType = "include"

	// FilterExclude excludes the specified nodes
	// from the aggregation.
	FilterExclude FilterType = "exclude"
)

// AggregationConfig holds the complete configuration required
// to perform an aggregation operation.
type AggregationConfig struct {
	// Aggregation defines the aggregation computation
	// and node-level behavior.
	Aggregation AggregationObject `json:"aggregation"`

	// Interval defines the time window over which
	// aggregation is performed.
	Interval IntervalObject `json:"interval"`

	// ResponseOptions controls formatting options
	// for the aggregation response.
	ResponseOptions ResponseOptions `json:"responseOptions,omitempty"`

	// Filter optionally limits which nodes are
	// included or excluded from aggregation.
	Filter *FilterObject `json:"filter,omitempty"`
}

// AggregationObject defines the aggregation operation
// and node grouping behavior.
type AggregationObject struct {
	// Compute specifies the aggregation function
	// (sum, avg, min, max, etc.).
	Compute AggregationType `json:"compute"`

	// ForEachNode determines whether aggregation
	// is performed per node or across all nodes combined.
	//
	// If true, aggregation is computed separately per node.
	// If false, data from all nodes is aggregated together.
	ForEachNode bool `json:"forEachNode,omitempty"`
}

// IntervalObject defines the time interval
// used for aggregation.
type IntervalObject struct {
	// Measure specifies the unit of time
	// (minute, hour, day, etc.).
	Measure IntervalMeasure `json:"measure"`

	// Interval specifies how many units of Measure
	// make up one aggregation bucket.
	//
	// Example: Measure=minute, Interval=5 → 5-minute buckets.
	Interval int `json:"interval"`
}

// ResponseOptions defines optional parameters
// that control how aggregation results are returned.
type ResponseOptions struct {
	// Timezone specifies the timezone to use
	// when generating timestamps.
	//
	// Defaults to "UTC" if not provided.
	Timezone string `json:"timezone,omitempty"`
}

// FilterObject defines node-level filtering rules
// applied before aggregation.
type FilterObject struct {
	// Nodes is the list of node IDs
	// to include or exclude.
	Nodes []string `json:"nodes,omitempty"`

	// Type determines whether the listed nodes
	// are included or excluded.
	Type FilterType `json:"type,omitempty"`
}

// AggregateDataPoint represents a single aggregated
// data point produced by an aggregation operation.
type AggregateDataPoint struct {
	// Timestamp represents the start of the aggregation
	// interval in Unix milliseconds.
	Timestamp int64 `json:"timestamp"`

	// Aggregate is the computed aggregated value
	// for the interval.
	Aggregate float64 `json:"aggregate"`
}
