package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Request struct
type GetSnapshotRequest struct {
	Timestamp int64    `json:"timestamp"`
	Variable  string   `json:"variable"`
	Nodes     []string `json:"nodes"`
}

// Response struct
type GetSnapshotResponse struct {
	Success    bool                   `json:"success"`
	Error      string                 `json:"error"`
	ReasonCode string                 `json:"reasonCode,omitempty"`
	Data       map[string]interface{} `json:"data"` // { "nodeID": { "timestamp": int64, "value": any } }
	Count      int                    `json:"count"`
}

// Validation Logic
func validateGetSnapshotRequest(req *GetSnapshotRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Variable == "" {
		return errors.New("variable is required")
	}
	if req.Timestamp <= 0 {
		return errors.New("timestamp must be greater than 0 (Unix milliseconds)")
	}
	if len(req.Nodes) == 0 {
		return errors.New("at least one node ID is required")
	}
	for i, node := range req.Nodes {
		if node == "" {
			return fmt.Errorf("node ID at index %d is empty", i)
		}
	}
	return nil
}

// GetSnapshot - Get variable value at (or nearest before) a specific timestamp
func (dm *DataManagement) GetSnapshot(ctx context.Context, req *GetSnapshotRequest) (*GetSnapshotResponse, error) {
	// Validation input
	if err := validateGetSnapshotRequest(req); err != nil {
		return nil, err
	}

	// Require authorization
	if dm.authToken == "" {
		return nil, errors.New("authorization token is missing")
	}

	// Build endpoint URL
	url := fmt.Sprintf("%s/v1/data/snapshot", dm.baseURL)

	// Marshal request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to crreate HTTP request: %w", err)
	}

	// Set required headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", dm.authToken)

	// Execute request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode response body
	var result GetSnapshotResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	// Check API-level success first
	if !result.Success {
		return &result, fmt.Errorf("API error: %s (reasonCode: %s)", result.Error, result.ReasonCode)
	}

	//  Verify HTTP status
	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	return &result, nil
}
