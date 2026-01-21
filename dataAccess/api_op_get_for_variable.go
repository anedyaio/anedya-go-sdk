// Package dataAccess provides APIs to retrieve and manage
// time-series data for nodes within the Anedya platform.
package dataAccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetDataRequest represents the payload used to fetch
// time-series data for a variable across one or more nodes.
type GetDataRequest struct {
	// Variable is the name of the variable whose data is requested.
	Variable string `json:"variable"`

	// Nodes is the list of node IDs for which data should be retrieved.
	Nodes []string `json:"nodes"`

	// From is the start timestamp of the query range (Unix milliseconds).
	From int64 `json:"from"`

	// To is the end timestamp of the query range (Unix milliseconds).
	To int64 `json:"to"`

	// Limit restricts the maximum number of data points returned.
	// If omitted, the server applies a default limit.
	Limit int `json:"limit,omitempty"`

	// Order specifies the sort order of returned data points.
	// Allowed values are "asc" or "desc".
	Order string `json:"order,omitempty"`
}

// GetDataResponse represents the response returned by
// the Get Data API.
type GetDataResponse struct {
	// Success indicates whether the request was processed successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message when Success is false.
	Error string `json:"error"`

	// ReasonCode is the machine-readable error code
	// used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`

	// Variable is the name of the requested variable.
	Variable string `json:"variable"`

	// Count represents the total number of data points returned.
	Count int `json:"count"`

	// Data maps node IDs to their corresponding data points.
	Data map[string][]DataPoint `json:"data"`
}

// GetData retrieves time-series data for a variable across one or more nodes
// from the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload and mandatory fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Get Data API.
//  4. Decode the API response into GetDataResponse.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle, cancellation, and deadlines.
//   - req: Pointer to GetDataRequest containing query parameters.
//
// Returns:
//   - (*GetDataResponse, nil) if the data is fetched successfully.
//   - (nil, error) for validation or client-side failures.
//   - (*GetDataResponse, error) when the API responds with an error.
func (dm *DataManagement) GetData(
	ctx context.Context,
	req *GetDataRequest,
) (*GetDataResponse, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get data request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// variable name must be provided
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

	// validate timestamp range
	if req.From <= 0 || req.To <= 0 || req.From > req.To {
		return nil, &errors.AnedyaError{
			Message: "invalid from/to timestamp range",
			Err:     errors.ErrInvalidTimeRange,
		}
	}

	// validate order field
	if req.Order != "" && req.Order != "asc" && req.Order != "desc" {
		return nil, &errors.AnedyaError{
			Message: "order must be asc or desc",
			Err:     errors.ErrInvalidOrder,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/data/getData", dm.baseURL)

	// convert request to JSON
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

	// send HTTP request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp GetDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetData response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle HTTP or API-level errors
	if resp.StatusCode != http.StatusOK || !apiResp.Success {
		return &apiResp, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &apiResp, nil
}
