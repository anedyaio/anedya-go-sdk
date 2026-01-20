package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// DeleteNodeRequest represents the payload sent to the Delete Node API.
// It contains the unique identifier of the node that needs to be deleted.
type DeleteNodeRequest struct {
	// NodeID is the unique identifier of the node to be deleted.
	NodeID string `json:"nodeid"`
}

// DeleteNodeResponse represents the response returned by the Delete Node API.
type DeleteNodeResponse struct {
	// Success indicates whether the node was deleted successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`
}

// DeleteNode deletes a node from the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory field (NodeID).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Delete Node API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into DeleteNodeResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to DeleteNodeRequest containing the node identifier.
//
// Returns:
//   - error: Returns nil if the operation succeeds, otherwise a sentinel
//     error or *errors.AnedyaError if validation, network, or API errors occur.
func (nm *NodeManagement) DeleteNode(
	ctx context.Context,
	req *DeleteNodeRequest,
) error {

	// Validate request object
	if req == nil {
		return &errors.AnedyaError{
			Message: "delete node request cannot be nil",
			Err:     errors.ErrDeleteNodeRequestNil,
		}
	}

	// Validate NodeID
	if req.NodeID == "" {
		return &errors.AnedyaError{
			Message: "node id is required to delete node",
			Err:     errors.ErrDeleteNodeIDRequired,
		}
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/delete", nm.baseURL)

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode DeleteNode request",
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
		return &errors.AnedyaError{
			Message: "failed to build DeleteNode request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute DeleteNode request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON
	var apiResp DeleteNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode DeleteNode response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	fmt.Printf("%+v\n", apiResp)

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// API-level error
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Delete successful
	return nil
}
