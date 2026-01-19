package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// ListChildNodesRequest represents the request payload sent to the List Child Nodes API.
// It contains the parent node ID and optional pagination parameters.
type ListChildNodesRequest struct {
	// ParentId is the unique identifier of the parent node.
	ParentId string `json:"parentId"`

	// Limit specifies the maximum number of child nodes to return in a single request.
	// If not provided or out of bounds, a default of 100 is applied.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of records to skip before returning results.
	// Defaults to 0 if negative.
	Offset int `json:"offset,omitempty"`
}

// ChildNode represents a single child node returned by the List Child Nodes API.
type ChildNode struct {
	// ChildId is the unique identifier of the child node.
	ChildId string `json:"childId"`

	// Alias is the alias assigned to the child node under the parent node.
	Alias string `json:"alias"`

	// CreatedAt indicates the timestamp when the child node was created.
	CreatedAt int64 `json:"createdAt"`
}

// ListChildNodesResponse represents the response returned by the List Child Nodes API.
type ListChildNodesResponse struct {
	// Success indicates whether the request was processed successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`

	// TotalCount indicates the total number of child nodes associated with the parent.
	TotalCount int `json:"totalCount"`

	// Count indicates the number of child nodes returned in this response.
	Count int `json:"count"`

	// Next indicates the offset value for the next page of results.
	Next int `json:"next"`

	// Data contains the list of child nodes returned in this response.
	Data []ChildNode `json:"data"`
}

// ListChildNodes retrieves the list of child nodes associated with a given parent node.
//
// This method performs the following operations:
//  1. Validates the request payload and ensures ParentId is provided.
//  2. Applies default pagination values if Limit or Offset are out of bounds.
//  3. Marshals the request payload into JSON.
//  4. Constructs an HTTP POST request to the List Child Nodes API endpoint.
//  5. Executes the HTTP request using the NodeManagement's HTTP client.
//  6. Decodes the API response into ListChildNodesResponse.
//  7. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to ListChildNodesRequest containing parent node ID and pagination info.
//
// Returns:
//   - *ListChildNodesResponse: Contains child nodes and pagination metadata on success.
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API errors occur.
func (nm *NodeManagement) ListChildNodes(
	ctx context.Context,
	req *ListChildNodesRequest,
) (*ListChildNodesResponse, error) {

	// Validate request object
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "list child nodes request cannot be nil",
			Err:     errors.ErrListChildNodesRequestNil,
		}
	}

	// Validate ParentId
	if req.ParentId == "" {
		return nil, &errors.AnedyaError{
			Message: "parent id is required to list child nodes",
			Err:     errors.ErrListChildNodesParentIDRequired,
		}
	}

	// Apply default and boundary values for pagination
	if req.Limit <= 0 || req.Limit > 1000 {
		req.Limit = 100
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/child/list", nm.baseURL)

	// Marshal request payload into JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode ListChildNodes request",
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
		return nil, &errors.AnedyaError{
			Message: "failed to build ListChildNodes request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute ListChildNodes request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON
	var apiResp ListChildNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode ListChildNodes response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// API-level error
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Success
	return &apiResp, nil
}
