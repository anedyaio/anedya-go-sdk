package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// RemoveChildNodeRequest represents the request payload sent to the Remove Child Node API.
// It contains the parent node ID and the child node ID to be detached.
type RemoveChildNodeRequest struct {
	// ParentId is the unique identifier of the parent node.
	ParentId string `json:"parentId"`

	// ChildNode is the unique identifier of the child node
	// to be removed from the parent.
	ChildNode string `json:"childNode"`
}

// RemoveChildNodeResponse represents the response returned by the Remove Child Node API.
type RemoveChildNodeResponse struct {
	// Success indicates whether the removal operation was successful.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// if the operation was unsuccessful.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`
}

// RemoveChildNode detaches a child node from its parent node in the Anedya platform.
//
// The request is provided using a *RemoveChildNodeRequest structure,
// which specifies both the parent node ID and the child node ID to remove.
//
// The method performs the following operations:
//  1. Validates the request payload and mandatory fields.
//  2. Marshals the request into JSON.
//  3. Constructs an HTTP POST request to the Remove Child Node API endpoint.
//  4. Executes the request using NodeManagement's HTTP client.
//  5. Decodes the API response.
//  6. Maps API or network errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - req: Pointer to RemoveChildNodeRequest containing parent and child node IDs.
//
// Returns:
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API-level errors occur.
func (nm *NodeManagement) RemoveChildNode(
	ctx context.Context,
	req *RemoveChildNodeRequest,
) error {

	// Validate request object
	if req == nil {
		return &errors.AnedyaError{
			Message: "remove child node request cannot be nil",
			Err:     errors.ErrRemoveChildNodeRequestNil,
		}
	}

	// Validate mandatory ParentId
	if req.ParentId == "" {
		return &errors.AnedyaError{
			Message: "parent id is required to remove child node",
			Err:     errors.ErrRemoveChildNodeParentIDRequired,
		}
	}

	// Validate mandatory ChildNode
	if req.ChildNode == "" {
		return &errors.AnedyaError{
			Message: "child node id is required to remove child node",
			Err:     errors.ErrRemoveChildNodeChildIDRequired,
		}
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/child/remove", nm.baseURL)

	// Marshal request payload into JSON
	body, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode RemoveChildNode request",
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
			Message: "failed to build RemoveChildNode request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute RemoveChildNode request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode API response JSON
	var apiResp RemoveChildNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode RemoveChildNode response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Check API-level success
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Success
	return nil
}
