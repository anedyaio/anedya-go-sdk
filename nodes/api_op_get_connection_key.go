package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetConnectionKeyRequest represents the payload sent to the Get Connection Key API.
// It contains the unique identifier of the node whose connection key is being requested.
type GetConnectionKeyRequest struct {
	// NodeID is the unique identifier of the node
	// whose connection key is being retrieved.
	NodeID string `json:"nodeid"`
}

// GetConnectionKeyResponse represents the response returned by the Get Connection Key API.
type GetConnectionKeyResponse struct {
	// Success indicates whether the request was successful.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`

	// ConnectionKey is the connection key associated with the node.
	ConnectionKey string `json:"connectionKey,omitempty"`
}

// GetConnectionKey retrieves the connection key associated with a node in the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory field (NodeID).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Get Connection Key API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into GetConnectionKeyResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to GetConnectionKeyRequest containing the node identifier.
//
// Returns:
//   - string: The connection key associated with the node on success.
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API errors occur.
func (nm *NodeManagement) GetConnectionKey(
	ctx context.Context,
	req *GetConnectionKeyRequest,
) (string, error) {

	// Validate request object is not nil
	if req == nil {
		return "", errors.ErrGetConnectionKeyRequestNil
	}

	// Validate NodeID is provided
	if req.NodeID == "" {
		return "", errors.ErrGetConnectionKeyNodeIDRequired
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/getConnectionKey", nm.baseURL)

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return "", &errors.AnedyaError{
			Message: "failed to encode GetConnectionKey request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Build HTTP POST request with context
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", &errors.AnedyaError{
			Message: "failed to build GetConnectionKey request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return "", &errors.AnedyaError{
			Message: "failed to execute GetConnectionKey request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON into GetConnectionKeyResponse
	var apiResp GetConnectionKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", &errors.AnedyaError{
			Message: "failed to decode GetConnectionKey response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Check HTTP status code for success
	if resp.StatusCode != http.StatusOK {
		return "", errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Check API-level success flag
	if !apiResp.Success {
		return "", errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Return the connection key
	return apiResp.ConnectionKey, nil
}
