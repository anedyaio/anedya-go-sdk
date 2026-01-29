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

// GetLatestDataRequest represents the payload used to fetch
// the most recent data point of a variable for one or more nodes.
type GetLatestDataRequest struct {
	// Nodes is the list of node IDs for which the latest data is requested.
	Nodes []string `json:"nodes"`

	// Variable is the name of the variable whose latest value is requested.
	Variable string `json:"variable"`
}

// getLatestDataAPIResponse represents the raw response
// returned by the Get Latest Data API.
type getLatestDataAPIResponse struct {
	common.BaseResponse

	// Data maps node IDs to their latest data point.
	Data map[string]DataPoint `json:"data"`

	// Count represents the number of nodes for which data was returned.
	Count int `json:"count"`
}

// GetLatestDataResult represents the processed and user-facing
// result returned by the GetLatestData method.
type GetLatestDataResult struct {
	// Data maps node IDs to their latest data point.
	Data map[string]DataPoint

	// Count represents the number of nodes for which data was returned.
	Count int
}

// GetLatestData retrieves the most recent data point of a variable
// for one or more nodes.
//
// Steps performed by this method:
//  1. Validate the request payload and required fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Get Latest Data API.
//  4. Decode the API response.
//  5. Convert the API response into a user-friendly result.
//
// Parameters:
//   - ctx: Context used to control request lifecycle, cancellation, and deadlines.
//   - req: Pointer to GetLatestDataRequest containing query parameters.
//
// Returns:
//   - (*GetLatestDataResult, nil) if the data is fetched successfully.
//   - (nil, error) for validation or client-side failures.
//   - (nil, error) when the API responds with an error.
func (dm *DataManagement) GetLatestData(
	ctx context.Context,
	req *GetLatestDataRequest,
) (*GetLatestDataResult, error) {

	// validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get latest data request cannot be nil",
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
	url := fmt.Sprintf("%s/v1/data/latest", dm.baseURL)

	// marshal request body
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

	// execute HTTP request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetLatestData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp getLatestDataAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetLatestData response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// return processed result
	return &GetLatestDataResult{
		Data:  apiResp.Data,
		Count: apiResp.Count,
	}, nil
}
