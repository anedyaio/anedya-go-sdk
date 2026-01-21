// Package dataAccess provides APIs to retrieve and manage
// time-series and latest variable data for nodes
// within the Anedya platform.
package dataAccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetLatestDataRequest represents the payload used to fetch
// the most recent data value for a variable across multiple nodes.
type GetLatestDataRequest struct {
	// Nodes is the list of node IDs for which the latest data is requested.
	Nodes []string `json:"nodes"`

	// Variable is the name of the variable whose latest value is requested.
	Variable string `json:"variable"`
}

// LatestDataPoint represents the latest recorded value
// of a variable for a node.
type LatestDataPoint struct {
	// Timestamp indicates when the value was last recorded (Unix milliseconds).
	Timestamp int64 `json:"timestamp"`

	// Value holds the most recent value of the variable.
	Value interface{} `json:"value"`
}

// GetLatestDataResponse represents the response returned by
// the Get Latest Data API.
type GetLatestDataResponse struct {
	// Success indicates whether the request was processed successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message when Success is false.
	Error string `json:"error"`

	// ReasonCode is the machine-readable error code
	// used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`

	// Data maps node IDs to their corresponding latest data points.
	Data map[string]LatestDataPoint `json:"data"`

	// Count represents the number of nodes for which data was returned.
	Count int `json:"count"`
}

// GetLatestData retrieves the most recent data value for a variable
// across one or more nodes from the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload and mandatory fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Get Latest Data API.
//  4. Decode the API response into GetLatestDataResponse.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle, cancellation, and deadlines.
//   - req: Pointer to GetLatestDataRequest containing query parameters.
//
// Returns:
//   - (*GetLatestDataResponse, nil) if the data is fetched successfully.
//   - (nil, error) for validation or client-side failures.
//   - (*GetLatestDataResponse, error) when the API responds with an error.
func (dm *DataManagement) GetLatestData(
	ctx context.Context,
	req *GetLatestDataRequest,
) (*GetLatestDataResponse, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get latest data request cannot be nil",
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

	// validate each node ID
	for i, node := range req.Nodes {
		if node == "" {
			return nil, &errors.AnedyaError{
				Message: fmt.Sprintf("node id at index %d is empty", i),
				Err:     errors.ErrInvalidNode,
			}
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/data/latest", dm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetLatestData request",
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
			Message: "failed to build GetLatestData request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// send HTTP request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetLatestData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp GetLatestDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetLatestData response",
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
