package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetConnectionKeyRequest defines the request payload used to
// fetch the current connection key for a node.
type GetConnectionKeyRequest struct {
	// NodeID represents the UUID of the node whose
	// connection key needs to be fetched.
	// This field is mandatory.
	NodeID string `json:"nodeid"`
}

// GetConnectionKeyResponse represents the API response structure
// returned by the server for both success and error cases.
type GetConnectionKeyResponse struct {
	// Success indicates whether the request was processed successfully.
	Success bool `json:"success"`

	// Error contains a descriptive error message if the request fails.
	Error string `json:"error"`

	// ReasonCode provides a machine-readable reason for failure (optional).
	ReasonCode string `json:"reasonCode,omitempty"`

	// ConnectionKey is the secret key used by the physical device
	// to authenticate and establish a connection.
	// Present only on successful responses.
	ConnectionKey string `json:"connectionKey,omitempty"`
}

// GetConnectionKey fetches the current connection key associated
// with a given node.
//
// The connection key is used by the physical device to authenticate
// itself and connect to the Anedya platform.
//
// The function returns an error if validation fails, the HTTP
// request fails, the API responds with an error, or the connection
// key is missing in a successful response.
func (nm *NodeManagement) GetConnectionKey(
	ctx context.Context,
	req *GetConnectionKeyRequest,
) (string, error) {

	// Validate request object
	if req == nil {
		return "", fmt.Errorf("request cannot be nil")
	}

	// Validate required field
	if req.NodeID == "" {
		return "", fmt.Errorf("nodeid is required and cannot be empty")
	}

	// Construct the API endpoint URL
	url := fmt.Sprintf("%s/v1/node/getConnectionKey", nm.baseURL)

	// Marshal request payload into JSON
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP POST request with context support
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required HTTP headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode API response (same structure for success and error)
	var keyResp GetConnectionKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&keyResp); err != nil {
		return "", fmt.Errorf("failed to decode API response: %w", err)
	}

	// Handle HTTP-level errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"API error %d: %s (reason: %s)",
			resp.StatusCode,
			keyResp.Error,
			keyResp.ReasonCode,
		)
	}

	// Handle application-level failure
	if !keyResp.Success {
		return "", fmt.Errorf(
			"failed to get connection key: %s (reason: %s)",
			keyResp.Error,
			keyResp.ReasonCode,
		)
	}

	// Final safety check — connection key must be present on success
	if keyResp.ConnectionKey == "" {
		return "", fmt.Errorf("connection key is empty in successful response")
	}

	// Return the fetched connection key
	return keyResp.ConnectionKey, nil
}
