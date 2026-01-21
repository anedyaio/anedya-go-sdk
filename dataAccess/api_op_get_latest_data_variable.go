package dataAccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// Request body for latest data API
type GetLatestDataRequest struct {
	Nodes    []string `json:"nodes"`    // list of node IDs
	Variable string   `json:"variable"` // variable name
}

// Single latest data value
type LatestDataPoint struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

// API response structure
type GetLatestDataResponse struct {
	Success    bool                       `json:"success"`
	Error      string                     `json:"error"`
	ReasonCode string                     `json:"reasonCode,omitempty"`
	Data       map[string]LatestDataPoint `json:"data"` // nodeID -> latest data
	Count      int                        `json:"count"`
}

// Fetches latest data for a variable
func (dm *DataManagement) GetLatestData(ctx context.Context, req *GetLatestDataRequest) (*GetLatestDataResponse, error) {

	// Validate request
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

	// Build API URL
	url := fmt.Sprintf("%s/v1/data/latest", dm.baseURL)

	// Convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetLatestData request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build GetLatestData request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Send request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetLatestData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Read response
	var apiResp GetLatestDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetLatestData response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Handle API error
	if resp.StatusCode != http.StatusOK || !apiResp.Success {
		return &apiResp, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return &apiResp, nil
}
