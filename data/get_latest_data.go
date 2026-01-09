package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Request: GetLatestDataRequest
type GetLatestDataRequest struct {
	Nodes    []string `json:"nodes"`
	Variable string   `json:"variable"`
}

// Response: GetLatestDataResponse
type GetLatestDataResponse struct {
	Success    bool                   `json:"success"`
	Error      string                 `json:"error"`
	ReasonCode string                 `json:"reasonCode,omitempty"`
	Data       map[string]interface{} `json:"data"`
	Count      int                    `json:"count"` // { "nodeID": { "timestamp": int64, "value": any } }
}

// Validation
func ValidateGetLatestDataRequest(req *GetLatestDataRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Variable == "" {
		return errors.New("variable is required")
	}
	if len(req.Nodes) == 0 {
		return errors.New("at least one node ID is required")
	}
	for i, node := range req.Nodes {
		if node == "" {
			return fmt.Errorf("node ID at %d index is empty", i)
		}
	}
	return nil
}

// GetLatestData - Fetch the most recent data point for a variable
func (dm *DataManagement) GetLatestData(ctx context.Context, req *GetLatestDataRequest) (*GetLatestDataResponse, error) {
	// Validate input
	if err := ValidateGetLatestDataRequest(req); err != nil {
		return nil, err
	}

	// Require auth token
	if dm.authToken == "" {
		return nil, errors.New("authorization token is missing")
	}

	// Endpoint URL
	url := fmt.Sprintf("%s/v1/data/latest", dm.baseURL)

	// Marshal request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", dm.authToken)

	// Execute request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	//Decode response
	var result GetLatestDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	// Check API-level success first
	if !result.Success {
		return &result, fmt.Errorf("API error: %s (reasonCode: %s)", result.Error, result.ReasonCode)
	}

	//Verify HTTP status
	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}
	return &result, nil
}
