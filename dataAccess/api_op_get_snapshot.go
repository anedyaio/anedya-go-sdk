package dataAccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// Request body for snapshot data
type GetSnapshotRequest struct {
	Timestamp int64    `json:"timestamp"` // time in unix ms
	Variable  string   `json:"variable"`  // variable name
	Nodes     []string `json:"nodes"`     // node IDs
}

// Single snapshot value
type SnapshotDataPoint struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

// API response structure
type GetSnapshotResponse struct {
	Success    bool                         `json:"success"`
	Error      string                       `json:"error"`
	ReasonCode string                       `json:"reasonCode,omitempty"`
	Data       map[string]SnapshotDataPoint `json:"data"` // nodeID -> snapshot
	Count      int                          `json:"count"`
}

// Fetches snapshot data for a variable at a given timestamp
func (dm *DataManagement) GetSnapshot(ctx context.Context, req *GetSnapshotRequest) (*GetSnapshotResponse, error) {

	// Validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get snapshot request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	if req.Variable == "" {
		return nil, &errors.AnedyaError{
			Message: "variable is required",
			Err:     errors.ErrVariableRequired,
		}
	}

	if req.Timestamp <= 0 {
		return nil, &errors.AnedyaError{
			Message: "timestamp must be greater than 0",
			Err:     errors.ErrInvalidTimestamp,
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
	url := fmt.Sprintf("%s/v1/data/snapshot", dm.baseURL)

	// Convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetSnapshot request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build GetSnapshot request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Send request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetSnapshot request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Read response
	var apiResp GetSnapshotResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetSnapshot response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Handle API error
	if resp.StatusCode != http.StatusOK || !apiResp.Success {
		return &apiResp, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return &apiResp, nil
}
