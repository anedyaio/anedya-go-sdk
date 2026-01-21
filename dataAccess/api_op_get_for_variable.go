// Package dataaccess provides APIs to fetch historical data
package dataAccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// DataManagement handles data-related APIs
type DataManagement struct {
	httpClient *http.Client
	baseURL    string
}

// Creates a new DataManagement client
func NewDataManagement(httpClient *http.Client, baseURL string) *DataManagement {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &DataManagement{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// Request body for GetData API
type GetDataRequest struct {
	Variable string   `json:"variable"`
	Nodes    []string `json:"nodes"`
	From     int64    `json:"from"` // start time (ms)
	To       int64    `json:"to"`   // end time (ms)
	Limit    int      `json:"limit,omitempty"`
	Order    string   `json:"order,omitempty"` // asc or desc
}

// Single data point
type DataPoint struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

// Response from GetData API
type GetDataResponse struct {
	Success    bool                   `json:"success"`
	Error      string                 `json:"error"`
	ReasonCode string                 `json:"reasonCode,omitempty"`
	Variable   string                 `json:"variable"`
	Count      int                    `json:"count"`
	Data       map[string][]DataPoint `json:"data"` // nodeId -> data points
}

// Fetches historical data for a variable
func (dm *DataManagement) GetData(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error) {

	// Basic validations
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get data request cannot be nil",
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

	if req.From <= 0 || req.To <= 0 || req.From > req.To {
		return nil, &errors.AnedyaError{
			Message: "invalid from/to timestamp range",
			Err:     errors.ErrInvalidTimeRange,
		}
	}

	if req.Order != "" && req.Order != "asc" && req.Order != "desc" {
		return nil, &errors.AnedyaError{
			Message: "order must be asc or desc",
			Err:     errors.ErrInvalidOrder,
		}
	}

	// API URL
	url := fmt.Sprintf("%s/v1/data/getData", dm.baseURL)

	// Convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetData request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build GetData request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Send request
	resp, err := dm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetData request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Read response
	var apiResp GetDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetData response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Handle API error
	if resp.StatusCode != http.StatusOK || !apiResp.Success {
		return &apiResp, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return &apiResp, nil
}
