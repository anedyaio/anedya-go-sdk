// Package nodes provides APIs to manage nodes and their hierarchical
// relationships within the Anedya platform.
package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// ChildNodeRequest represents a child node that is to be attached
// to a parent node along with a logical alias.
type ChildNodeRequest struct {
	// NodeId is the unique identifier of the child node.
	NodeId string `json:"nodeId"`

	// Alias is the logical alias assigned to the child node
	// under the parent node.
	Alias string `json:"alias"`
}

// AddChildNodeRequest represents the payload sent to the Add Child Node API.
// It contains the parent node identifier and a list of child nodes to attach.
type AddChildNodeRequest struct {
	// ParentId is the unique identifier of the parent node to which child nodes will be attached.
	ParentId string `json:"parentId"`

	// ChildNodes is the list of child nodes to be attached under the parent node.
	ChildNodes []ChildNodeRequest `json:"childNodes"`
}

// AddChildNodeResponse represents the response returned by the Add Child Node API.
type AddChildNodeResponse struct {
	// Success indicates whether the Add Child Node operation was successful.
	Success bool `json:"success"`

	// Error contains the human-readable error message when Success is false.
	Error string `json:"error"`

	// ReasonCode is the machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`
}

// AddChildNode attaches one or more child nodes to a parent node in the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload and mandatory fields.
//  2. Marshal the request into JSON.
//  3. Build and send a POST request to the Add Child Node API.
//  4. Decode the API response into AddChildNodeResponse.
//  5. Map API errors into structured SDK errors.
//
// Parameters:
//   - ctx: The context for controlling request cancellation and deadlines.
//   - req: Pointer to AddChildNodeRequest containing parent and child node details.
//
// Returns:
//   - error: Returns nil if the operation succeeds, or a structured error
//     (sentinel error or *errors.AnedyaError) if validation, network,
//     or API errors occur.
func (nm *NodeManagement) AddChildNode(ctx context.Context, req *AddChildNodeRequest) error {

	// Validate that request object is not nil
	if req == nil {
		return errors.ErrAddChildNodeRequestNil
	}

	// Validate that ParentId is provided
	if req.ParentId == "" {
		return errors.ErrAddChildNodeParentIdRequired
	}

	// Validate that at least one child node is provided
	if len(req.ChildNodes) == 0 {
		return errors.ErrAddChildNodeEmptyChildren
	}

	// Validate each child node: both NodeId and Alias are required
	for i, c := range req.ChildNodes {
		if c.NodeId == "" || c.Alias == "" {
			return &errors.AnedyaError{
				Message: fmt.Sprintf(
					"childNodes[%d] requires both nodeId and alias",
					i,
				),
				Err: fmt.Errorf("invalid child node entry"),
			}
		}
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/child/add", nm.baseURL)

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode AddChildNode request",
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
			Message: "failed to build AddChildNode request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute AddChildNode request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON into AddChildNodeResponse
	var apiResp AddChildNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode AddChildNode response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Check HTTP status code for success
	if resp.StatusCode != http.StatusOK {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Check API-level success flag
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Operation successful
	return nil
}
