// Package dataAccess provides APIs to retrieve and manage
// historical and latest time-series data for nodes
// within the Anedya platform.
package dataAccess

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

	// 1. Validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get latest data request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	if req.Variable == "" {
		return nil, &errors.AnedyaError{
			Message: "variable is required",
			Err:     errors.ErrVariableRequired,
		}
	}

	if len(req.Nodes) == 0 {
		return nil, &errors.AnedyaError{
			Message: "at least one node must be provided",
			Err:     errors.ErrNodesEmpty,
		}
	}

	for i, node := range req.Nodes {
		if node == "" {
			return nil, &errors.AnedyaError{
				Message: fmt.Sprintf("node id at index %d is empty", i),
				Err:     errors.ErrInvalidNode,
			}
		}
	}

	// 2. Encode request
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetLatestData request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/data/latest", dm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build GetLatestData request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetLatestData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read GetLatestData response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp getLatestDataAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetLatestData response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Return processed result
	return &GetLatestDataResult{
		Data:  apiResp.Data,
		Count: apiResp.Count,
	}, nil
}
