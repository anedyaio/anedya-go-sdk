// Package dataAccess provides APIs to retrieve and manage
// historical and latest time-series data for nodes
// within the Anedya platform.
package dataAccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// ========================
// REQUEST
// ========================

// GetDataRequest represents the payload used to fetch
// historical data points for a variable across one or more nodes
// within a given time range.
type GetDataRequest struct {
	// Variable is the name of the variable whose data is requested.
	Variable string `json:"variable"`

	// Nodes is the list of node IDs for which data is requested.
	Nodes []string `json:"nodes"`

	// From is the start timestamp (inclusive) in milliseconds.
	From int64 `json:"from"`

	// To is the end timestamp (inclusive) in milliseconds.
	To int64 `json:"to"`

	// Limit specifies the maximum number of data points to return per node.
	Limit int `json:"limit,omitempty"`

	// Order defines the sorting order of data points.
	// Allowed values: "asc" or "desc".
	Order string `json:"order,omitempty"`
}

// ========================
// INTERNAL API RESPONSE
// ========================

// getDataAPIResponse represents the raw response returned
// by the Get Data API.
type getDataAPIResponse struct {
	common.BaseResponse

	// Variable is the variable name for which data was fetched.
	Variable string `json:"variable"`

	// Count represents the total number of data points returned.
	Count int `json:"count"`

	// Data maps node IDs to their corresponding data points.
	Data map[string][]DataPoint `json:"data"`
}

// ========================
// RESULT
// ========================

// GetDataResult represents the processed and user-facing
// result returned by the GetData method.
type GetDataResult struct {
	// Variable is the variable name for which data was fetched.
	Variable string

	// Count represents the total number of data points returned.
	Count int

	// Data maps node IDs to their corresponding data points.
	Data map[string][]DataPoint
}

// ========================
// API
// ========================

// GetData retrieves historical data points for a variable
// across one or more nodes within a specified time range.
//
// Steps performed by this method:
//  1. Validate the request payload and required fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Get Data API.
//  4. Decode the API response.
//  5. Convert the API response into a user-friendly result.
//
// Parameters:
//   - ctx: Context used to control request lifecycle, cancellation, and deadlines.
//   - req: Pointer to GetDataRequest containing query parameters.
//
// Returns:
//   - (*GetDataResult, nil) if the data is fetched successfully.
//   - (nil, error) for validation or client-side failures.
//   - (nil, error) when the API responds with an error.
func (dm *DataManagement) GetData(
	ctx context.Context,
	req *GetDataRequest,
) (*GetDataResult, error) {

	// validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get data request cannot be nil",
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

	// at least one node must be provided
	if len(req.Nodes) == 0 {
		return nil, &errors.AnedyaError{
			Message: "at least one node must be provided",
			Err:     errors.ErrNodesEmpty,
		}
	}

	// validate time range
	if req.From <= 0 || req.To <= 0 || req.From > req.To {
		return nil, &errors.AnedyaError{
			Message: "invalid from/to timestamp range",
			Err:     errors.ErrInvalidTimeRange,
		}
	}

	// validate order value
	if req.Order != "" && req.Order != "asc" && req.Order != "desc" {
		return nil, &errors.AnedyaError{
			Message: "order must be asc or desc",
			Err:     errors.ErrInvalidOrder,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/data/getData", dm.baseURL)

	// marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetData request",
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
			Message: "failed to build GetData request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// execute HTTP request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp getDataAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetData response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle HTTP-level errors
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// return processed result
	return &GetDataResult{
		Variable: apiResp.Variable,
		Count:    apiResp.Count,
		Data:     apiResp.Data,
	}, nil
}
