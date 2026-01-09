package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ClearChildNodesRequest defines the request payload used to
// remove all child nodes associated with a given parent node.
type ClearChildNodesRequest struct {
	// ParentId represents the UUID of the parent node
	// whose child nodes need to be cleared.
	// This field is mandatory.
	ParentId string `json:"parentId"`
}

// ClearChildNodesResponse represents the API response structure
// returned by the server for both success and error cases.
type ClearChildNodesResponse struct {
	// Success indicates whether the operation completed successfully.
	Success bool `json:"success"`

	// Error contains a descriptive error message if the operation fails.
	Error string `json:"error"`

	// ReasonCode provides a machine-readable reason for failure (optional).
	ReasonCode string `json:"reasonCode,omitempty"`
}

// ClearAllChildNodes removes ALL child node associations linked
// to the specified parent node.
//
// ⚠️ This is a destructive operation:
// All existing child relationships will be permanently removed.
// Use this method with caution.
//
// The function returns an error if validation fails, the HTTP
// request fails, or the API responds with an error.
func (nm *NodeManagement) ClearAllChildNodes(
	ctx context.Context,
	req *ClearChildNodesRequest,
) error {

	// Validate request object
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate required field
	if req.ParentId == "" {
		return fmt.Errorf("parentId is required and cannot be empty")
	}

	// Construct the API endpoint URL
	url := fmt.Sprintf("%s/v1/node/child/clear", nm.baseURL)

	// Marshal request payload into JSON
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP POST request with context support
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required HTTP headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode API response (same structure for success and error)
	var clearResp ClearChildNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&clearResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// Handle HTTP-level errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"API error %d: %s (reason: %s)",
			resp.StatusCode,
			clearResp.Error,
			clearResp.ReasonCode,
		)
	}

	// Handle application-level failure
	if !clearResp.Success {
		return fmt.Errorf(
			"failed to clear child nodes: %s (reason: %s)",
			clearResp.Error,
			clearResp.ReasonCode,
		)
	}

	// Child nodes successfully cleared
	return nil
}
