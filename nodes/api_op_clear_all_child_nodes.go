package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// ClearChildNodesRequest represents the payload sent to the Clear Child Nodes API.
// It contains the parent node identifier whose child nodes are to be removed.
type ClearChildNodesRequest struct {
	// ParentId is the unique identifier of the parent node
	// whose child nodes will be removed.
	ParentId string `json:"parentId"`
}

// ClearChildNodesResponse represents the response returned by the Clear Child Nodes API.
type ClearChildNodesResponse struct {
	// Success indicates whether the operation was completed successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`
}

// ClearChildNodes removes all child nodes associated with a given parent node in the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory fields (ParentId).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Clear Child Nodes API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into ClearChildNodesResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to ClearChildNodesRequest containing the parent node identifier.
//
// Returns:
//   - error: Returns nil if the operation succeeds, otherwise returns a sentinel
//     error or *errors.AnedyaError if validation, network, or API errors occur.
func (nm *NodeManagement) ClearChildNodes(
	ctx context.Context,
	req *ClearChildNodesRequest,
) error {

	// Validate request object
	if req == nil {
		return &errors.AnedyaError{
			Message: "clear child nodes request cannot be nil",
			Err:     errors.ErrClearChildNodesRequestNil,
		}
	}

	// Validate ParentId
	if req.ParentId == "" {
		return &errors.AnedyaError{
			Message: "parentId is required to clear child nodes",
			Err:     errors.ErrClearChildNodesParentIDRequired,
		}
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/child/clear", nm.baseURL)

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode ClearChildNodes request",
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
			Message: "failed to build ClearChildNodes request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute ClearChildNodes request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON
	var apiResp ClearChildNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode ClearChildNodes response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle HTTP or API level error
	if resp.StatusCode != http.StatusOK || !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
