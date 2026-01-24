// Package aggregations provides APIs to compute and retrieve
// aggregated views of variable data over time within the
// Anedya platform.
package aggregations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetAggregationByTimeRequest represents the request payload for fetching
// time-based aggregations of variable data.
//
// Time aggregations are useful for answering questions such as:
//   - What is the average value of a variable over a given time range?
//   - What is the maximum value in every 15-minute interval?
//   - What is the difference between the first and last value in a week?
type GetAggregationByTimeRequest struct {
	// Variable is the name of the variable to aggregate.
	Variable string `json:"variable"`

	// From is the start timestamp (inclusive) in Unix milliseconds.
	From int64 `json:"from"`

	// To is the end timestamp (inclusive) in Unix milliseconds.
	To int64 `json:"to"`

	// Config defines how the aggregation should be computed,
	// including aggregation type, interval, and filters.
	Config AggregationConfig `json:"config"`
}

// getAggregationAPIResponse represents the raw response returned
// by the Get Aggregation By Time API.
type getAggregationAPIResponse struct {
	common.BaseResponse

	// Variable is the variable name for which aggregation was performed.
	Variable string `json:"variable"`

	// Config is the aggregation configuration used for computation.
	Config AggregationConfig `json:"config"`

	// Data maps node IDs to their corresponding aggregated data points.
	//
	// Behavior depends on the aggregation configuration:
	//   - If forEachNode=true: each node ID has its own aggregated series.
	//   - If forEachNode=false: aggregation may be grouped under a single key.
	Data map[string][]AggregateDataPoint `json:"data"`
}

// GetAggregationResult represents the processed and user-facing
// result returned by the GetAggregationByTime method.
type GetAggregationResult struct {
	// Variable is the variable name for which aggregation was performed.
	Variable string

	// Config is the aggregation configuration used.
	Config AggregationConfig

	// Data maps node IDs to their corresponding aggregated data points.
	//
	// Key structure:
	//   - forEachNode=true  → each node ID has its own aggregation
	//   - forEachNode=false → typically a single aggregated result
	Data map[string][]AggregateDataPoint
}

// GetAggregationByTime retrieves time-based aggregations for a variable
// across nodes within a specified time range.
//
// This API is designed to efficiently summarize large volumes of
// time-series data into meaningful aggregated results. Supported
// aggregation types include:
//
//   - sum       : Sum of values over each interval
//   - avg       : Average value
//   - median    : Median value
//   - min / max : Minimum or maximum value
//   - diff      : Difference between first and last value
//   - deltasum  : Difference accounting for counter resets
//   - stddev    : Standard deviation
//
// Notes:
//   - Aggregation APIs incur charges based on GB scanned.
//   - Maximum processing time per request is 10 seconds.
//
// Steps performed by this method:
//  1. Validate the request payload and required fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the aggregation API.
//  4. Decode the API response.
//  5. Convert the API response into a user-friendly result.
//
// Parameters:
//   - ctx: Context used to control request lifecycle, cancellation, and deadlines.
//   - req: Pointer to GetAggregationByTimeRequest containing query parameters.
//
// Returns:
//   - (*GetAggregationResult, nil) if aggregation is computed successfully.
//   - (nil, error) for validation or client-side failures.
//   - (nil, error) when the API responds with an error.
func (ac *AggregationsClient) GetAggregationByTime(
	ctx context.Context,
	req *GetAggregationByTimeRequest,
) (*GetAggregationResult, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "aggregation request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// variable name is mandatory
	if req.Variable == "" {
		return nil, &errors.AnedyaError{
			Message: "variable is required",
			Err:     errors.ErrVariableRequired,
		}
	}

	// validate time range
	if req.From <= 0 || req.To <= 0 || req.From > req.To {
		return nil, &errors.AnedyaError{
			Message: "invalid from/to timestamp range",
			Err:     errors.ErrInvalidTimeRange,
		}
	}

	// aggregation compute type must be provided
	if req.Config.Aggregation.Compute == "" {
		return nil, &errors.AnedyaError{
			Message: "aggregation compute type is required",
			Err:     errors.ErrInvalidInput,
		}
	}

	// validate interval configuration
	if req.Config.Interval.Measure == "" || req.Config.Interval.Interval <= 0 {
		return nil, &errors.AnedyaError{
			Message: "valid interval measure and interval value are required",
			Err:     errors.ErrInvalidInput,
		}
	}

	// validate filter configuration if provided
	if req.Config.Filter != nil {
		if len(req.Config.Filter.Nodes) == 0 {
			return nil, &errors.AnedyaError{
				Message: "filter nodes cannot be empty when filter is provided",
				Err:     errors.ErrInvalidInput,
			}
		}
		if req.Config.Filter.Type != FilterInclude && req.Config.Filter.Type != FilterExclude {
			return nil, &errors.AnedyaError{
				Message: "filter type must be 'include' or 'exclude'",
				Err:     errors.ErrInvalidInput,
			}
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/aggregates/variable/byTime", ac.baseURL)

	// marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode aggregation request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// create HTTP request with context
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build aggregation request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// execute HTTP request
	resp, err := ac.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute aggregation request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp getAggregationAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode aggregation response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Check for any error (HTTP or API-level)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices || !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// return processed result
	return &GetAggregationResult{
		Variable: apiResp.Variable,
		Config:   apiResp.Config,
		Data:     apiResp.Data,
	}, nil
}
