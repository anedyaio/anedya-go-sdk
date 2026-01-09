package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Data Management
type DataManagement struct {
	httpClient *http.Client
	baseURL    string
	authToken  string // Full Authorization header value eg:"Bearer xxxx"
}

// NewDataManagement creates a new Client instance.
// authToken should include the prefix is needed
func NewDataManagement(httpClient *http.Client, baseURL string, authToken string) *DataManagement {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &DataManagement{
		httpClient: httpClient,
		baseURL:    baseURL,
		authToken:  authToken,
	}
}

// Request Struct
type GetDataRequest struct {
	Variable string   `json:"variable"`
	Nodes    []string `json:"nodes"`
	From     int64    `json:"from"` // Unix milliseconds
	To       int64    `json:"to"`   // Unix milliseconds
	Limit    int      `json:"limit,omitempty"`
	Order    string   `json:"order,omitempty"` // asc | desc
}

// Response Struct
type GetDataResponse struct {
	Success    bool                   `json:"success"`
	Error      string                 `json:"error"`
	ReasonCode string                 `json:"reasonCode,omitempty"`
	Variable   string                 `json:"variable"`
	Count      int                    `json:"count"`
	Data       map[string]interface{} `json:"data"` // { "nodeID": [ { "timestamp": int64, "value": any } ... ] }
}

// Validation Logic
func ValidateGetDataRequest(req *GetDataRequest) error {
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
			return fmt.Errorf("node ID at index %d is empty", i)
		}
	}
	if req.From <= 0 {
		return errors.New("from timestamp must be greater than 0")
	}
	if req.To <= 0 {
		return errors.New("to timestamp must be greater than 0")
	}
	if req.From > req.To {
		return errors.New("from cannot be greater than to")
	}
	if req.Limit < 0 {
		return errors.New("limit cannot be negative")
	}
	if req.Order != "" && req.Order != "asc" && req.Order != "desc" {
		return errors.New("order must be 'asc' or 'desc' if provided")
	}
	return nil
}

// GetData - Fetch historical data
func (dm *DataManagement) GetData(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error) {
	// Input validation
	if err := ValidateGetDataRequest(req); err != nil {
		return nil, err
	}

	// Check for auth token
	if dm.authToken == "" {
		return nil, errors.New("authorization token is missing")
	}

	// Build endpoint URL
	url := fmt.Sprintf("%s/v1/data/getData", dm.baseURL)

	// Marshal request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", dm.authToken)

	// Execute Request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Always try to Decode the JSON response Body
	var result GetDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	// First: Check API-level success
	if !result.Success {
		return &result, fmt.Errorf("API error: %s (reasonCode: %s)", result.Error, result.ReasonCode)
	}

	// Then: Verify HTTP status
	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	return &result, err
}
