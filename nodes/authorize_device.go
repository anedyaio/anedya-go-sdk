package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthorizeDeviceRequest defines the request payload used to
// authorize (bind) a physical device to a node.
//
// This request is sent from the server side and skips the
// normal device-initiated provisioning flow.
type AuthorizeDeviceRequest struct {
	// NodeID represents the target node UUID.
	// This field is mandatory.
	NodeID string `json:"nodeId"`

	// DeviceID represents the unique identifier of the physical device
	// (e.g., MAC address, serial number).
	// This field is mandatory.
	DeviceID string `json:"deviceId"`
}

// AuthorizeDeviceResponse represents the API response structure
// returned by the server for both success and error cases.
type AuthorizeDeviceResponse struct {
	// Success indicates whether the authorization was successful.
	Success bool `json:"success"`

	// Error contains a descriptive error message if the request fails.
	Error string `json:"error"`

	// ReasonCode provides a machine-readable reason for failure (optional).
	ReasonCode string `json:"reasonCode,omitempty"`
}

// AuthorizeDevice directly binds a physical device to a node
// using a server-side API call.
//
// This method bypasses the device-side provisioning process and
// is useful for pre-provisioning, factory setup, testing,
// or recovery scenarios.
//
// The function returns an error if validation fails, the HTTP
// request fails, or the API responds with an error.
func (nm *NodeManagement) AuthorizeDevice(
	ctx context.Context,
	req *AuthorizeDeviceRequest,
) error {

	// Validate request object
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate required fields
	if req.NodeID == "" {
		return fmt.Errorf("nodeid is required and cannot be empty")
	}
	if req.DeviceID == "" {
		return fmt.Errorf("deviceid is required and cannot be empty")
	}

	// Construct the API endpoint URL
	url := fmt.Sprintf("%s/v1/node/authorize", nm.baseURL)

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
	var authResp AuthorizeDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// Handle HTTP-level errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"API error %d: %s (reason: %s)",
			resp.StatusCode,
			authResp.Error,
			authResp.ReasonCode,
		)
	}

	// Handle application-level failure
	if !authResp.Success {
		return fmt.Errorf(
			"device authorization failed: %s (reason: %s)",
			authResp.Error,
			authResp.ReasonCode,
		)
	}

	// Authorization successful
	return nil
}
