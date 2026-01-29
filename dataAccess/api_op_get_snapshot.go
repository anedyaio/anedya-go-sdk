// Package dataAccess provides APIs to retrieve historical,
// latest, and snapshot data for nodes within the Anedya platform.
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

// GetSnapshotRequest represents the payload used to fetch
// snapshot data of a variable for one or more nodes at
// a specific timestamp.
type GetSnapshotRequest struct {
	// Timestamp is the point in time (Unix milliseconds)
	// at which the snapshot data is requested.
	Timestamp int64 `json:"timestamp"`

	// Variable is the name of the variable whose snapshot is requested.
	Variable string `json:"variable"`

	// Nodes is the list of node IDs for which snapshot data is requested.
	Nodes []string `json:"nodes"`
}

// getSnapshotAPIResponse represents the raw API response
// returned by the Snapshot Data API.
type getSnapshotAPIResponse struct {
	common.BaseResponse

	// Data maps node IDs to their snapshot data point.
	Data map[string]DataPoint `json:"data"`

	// Count represents the number of nodes for which data was returned.
	Count int `json:"count"`
}

// GetSnapshotResult represents the processed and user-facing
// result returned by the GetSnapshot method.
type GetSnapshotResult struct {
	// Data maps node IDs to their snapshot data point.
	Data map[string]DataPoint

	// Count represents the number of nodes for which data was returned.
	Count int
}

// GetSnapshot retrieves snapshot data of a variable for one or more
// nodes at a specific timestamp.
//
// Steps performed by this method:
//  1. Validate request parameters.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Snapshot Data API.
//  4. Decode the API response.
//  5. Convert the API response into a user-friendly result.
//
// Parameters:
//   - ctx: Context used to control request lifecycle, cancellation, and deadlines.
//   - req: Pointer to GetSnapshotRequest containing snapshot query parameters.
//
// Returns:
//   - (*GetSnapshotResult, nil) if snapshot data is fetched successfully.
//   - (nil, error) for validation or client-side failures.
//   - (nil, error) when the API responds with an error.
func (dm *DataManagement) GetSnapshot(
	ctx context.Context,
	req *GetSnapshotRequest,
) (*GetSnapshotResult, error) {

	// validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get snapshot request cannot be nil",
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

	// timestamp must be valid
	if req.Timestamp <= 0 {
		return nil, &errors.AnedyaError{
			Message: "timestamp must be greater than 0",
			Err:     errors.ErrInvalidTimestamp,
		}
	}

	// at least one node must be provided
	if len(req.Nodes) == 0 {
		return nil, &errors.AnedyaError{
			Message: "at least one node must be provided",
			Err:     errors.ErrNodesEmpty,
		}
	}

	// validate node IDs
	for i, node := range req.Nodes {
		if node == "" {
			return nil, &errors.AnedyaError{
				Message: fmt.Sprintf("node id at index %d is empty", i),
				Err:     errors.ErrInvalidNode,
			}
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/data/snapshot", dm.baseURL)

	// marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetSnapshot request",
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
			Message: "failed to build GetSnapshot request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// execute HTTP request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetSnapshot request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp getSnapshotAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetSnapshot response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// return processed result
	return &GetSnapshotResult{
		Data:  apiResp.Data,
		Count: apiResp.Count,
	}, nil
}
