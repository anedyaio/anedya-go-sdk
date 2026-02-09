// Package aggregations provides APIs to compute and retrieve
// aggregated views of variable data over time within the
// Anedya platform.
package aggregations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// AggregationType represents the type of aggregation
// operation to be performed on time-series data.
type AggregationType string

const (
	// AggregationSum computes the sum of values
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

	// AggregationDeltaSum computes the delta sum while
	// accounting for counter resets.
	AggregationDeltaSum AggregationType = "deltasum"

	// AggregationStdDev computes the standard deviation
	// over the specified interval.
	AggregationStdDev AggregationType = "stddev"
)

// IntervalMeasure represents the unit of time used
// for aggregation intervals.
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
// required to perform a time-based aggregation.
type AggregationConfig struct {
	// Aggregation defines the aggregation method and options.
	Aggregation AggregationObject `json:"aggregation"`

	// Interval defines the time-based aggregation window.
	Interval IntervalObject `json:"interval"`

	// ResponseOptions defines formatting options
	// for the aggregation response.
	ResponseOptions ResponseOptions `json:"responseOptions,omitempty"`

	// Filter defines optional node filtering rules.
	Filter *FilterObject `json:"filter,omitempty"`
}

// AggregationObject defines the aggregation operation
// to be performed.
type AggregationObject struct {
	// Compute specifies the aggregation type
	// (sum, avg, min, max, etc.).
	Compute AggregationType `json:"compute"`

	// ForEachNode determines whether aggregation is computed
	// per node or across all nodes.
	ForEachNode bool `json:"forEachNode,omitempty"`
}

// IntervalObject defines the size and unit of the
// aggregation time window.
type IntervalObject struct {
	// Measure specifies the unit of time
	// (minute, hour, day, etc.).
	Measure IntervalMeasure `json:"measure"`

	// Interval specifies the number of units
	// for the aggregation window.
	Interval int `json:"interval"`
}

// ResponseOptions defines options for formatting
// aggregation results.
type ResponseOptions struct {
	// Timezone specifies the timezone used for
	// response timestamps (default: UTC).
	Timezone string `json:"timezone,omitempty"`
}

// FilterObject defines node-level filtering rules
// applied before aggregation.
type FilterObject struct {
	// Nodes is the list of node IDs used for filtering.
	Nodes []string `json:"nodes,omitempty"`

	// Type specifies whether nodes are included or excluded.
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

// GetAggregationByTimeRequest represents the payload used
// to request aggregated data for a variable over time.
type GetAggregationByTimeRequest struct {
	// Variable is the name of the variable to aggregate.
	Variable string `json:"variable"`

	// From is the start timestamp (Unix milliseconds).
	From int64 `json:"from"`

	// To is the end timestamp (Unix milliseconds).
	To int64 `json:"to"`

	// Config defines the aggregation configuration.
	Config AggregationConfig `json:"config"`
}

// getAggregationAPIResponse represents the raw API response
// returned by the aggregation endpoint.
type getAggregationAPIResponse struct {
	common.BaseResponse
	Variable string                          `json:"variable"`
	Config   AggregationConfig               `json:"config"`
	Data     map[string][]AggregateDataPoint `json:"data"`
}

// GetAggregationResult represents the processed result
// returned to SDK consumers.
type GetAggregationResult struct {
	// Variable is the name of the aggregated variable.
	Variable string

	// Config is the aggregation configuration used.
	Config AggregationConfig

	// Data maps node IDs to their aggregated data points.
	Data map[string][]AggregateDataPoint
}

// GetAggregationByTime retrieves aggregated data for a variable
// over a specified time range.
//
// Steps performed by this method:
//  1. Validate request payload and aggregation configuration.
//  2. Marshal the request into JSON.
//  3. Send a POST request to the aggregation API.
//  4. Decode and validate the API response.
//  5. Map API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle and cancellation.
//   - req: Pointer to GetAggregationByTimeRequest containing query details.
//
// Returns:
//   - (*GetAggregationResult, nil) on successful aggregation.
//   - (nil, error) for validation, network, or API errors.
func (ac *AggregationsManagement) GetAggregationByTime(
	ctx context.Context,
	req *GetAggregationByTimeRequest,
) (*GetAggregationResult, error) {

	// 1. Validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "aggregation request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	if req.Variable == "" {
		return nil, &errors.AnedyaError{
			Message: "variable is required",
			Err:     errors.ErrVariableRequired,
		}
	}

	if req.From <= 0 || req.To <= 0 || req.From > req.To {
		return nil, &errors.AnedyaError{
			Message: "invalid from/to timestamp range",
			Err:     errors.ErrInvalidTimeRange,
		}
	}

	if req.Config.Aggregation.Compute == "" {
		return nil, &errors.AnedyaError{
			Message: "aggregation compute type is required",
			Err:     errors.ErrAggregationComputeRequired,
		}
	}

	if req.Config.Interval.Measure == "" || req.Config.Interval.Interval <= 0 {
		return nil, &errors.AnedyaError{
			Message: "invalid aggregation interval configuration",
			Err:     errors.ErrInvalidInterval,
		}
	}

	if req.Config.Filter != nil {
		if len(req.Config.Filter.Nodes) == 0 {
			return nil, &errors.AnedyaError{
				Message: "filter nodes cannot be empty",
				Err:     errors.ErrFilterNodesRequired,
			}
		}

		if req.Config.Filter.Type != FilterInclude && req.Config.Filter.Type != FilterExclude {
			return nil, &errors.AnedyaError{
				Message: "filter type must be include or exclude",
				Err:     errors.ErrInvalidFilterType,
			}
		}
	}

	// 2. Build API URL
	url := fmt.Sprintf("%s/v1/aggregates/variable/byTime", ac.baseURL)

	// 3. Marshal request
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode aggregation request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 4. Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build aggregation request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 5. Execute HTTP request
	resp, err := ac.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute aggregation request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 6. Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read aggregation response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 7. Decode response
	var apiResp getAggregationAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode aggregation response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 8. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 9. Return result
	return &GetAggregationResult{
		Variable: apiResp.Variable,
		Config:   apiResp.Config,
		Data:     apiResp.Data,
	}, nil
}
