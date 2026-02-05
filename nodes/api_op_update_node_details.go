package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// UpdateType represents the type of update operation
// that can be performed on a node.
type UpdateType string

const (
	// UpdateNodeName updates the node's name
	UpdateNodeName UpdateType = "node_name"

	// UpdateNodeDesc updates the node's description
	UpdateNodeDesc UpdateType = "node_desc"

	// UpdateTag adds, updates, or removes a tag
	UpdateTag UpdateType = "tag"
)

// NodeUpdate represents a single update operation
// applied to a node.
type NodeUpdate struct {

	// Type specifies the kind of update to perform.
	// Supported values are:
	//   - node_name
	//   - node_desc
	//   - tag
	Type UpdateType `json:"type"`

	// Value contains the new value for name or description updates.
	// This field is mandatory for non-tag updates.
	Value string `json:"value,omitempty"`

	// Tag contains the tag object for tag-related updates.
	// This field is mandatory when Type is UpdateTag.
	Tag *Tag `json:"tag,omitempty"`
}

// UpdateNodeRequest represents the payload sent to
// the Update Node API.
type UpdateNodeRequest struct {

	// NodeID is the unique identifier of the node
	// to be updated.
	NodeID string `json:"nodeid"`

	// Updates contains one or more update operations
	// to be applied to the node.
	Updates []NodeUpdate `json:"updates"`
}

// UpdateNodeResponse represents the response returned
// by the Update Node API.
type UpdateNodeResponse struct {
	common.BaseResponse
}

// UpdateNode applies one or more updates to a node.
//
// The method performs the following steps:
//  1. Validates the request payload and update operations.
//  2. Marshals the request into JSON.
//  3. Sends the request to the Update Node API.
//  4. Decodes and validates the API response.
//
// Validation failures return sentinel errors defined in
// the errors package. Network and API failures return
// *errors.AnedyaError.
func (nm *NodeManagement) UpdateNode(
	ctx context.Context,
	req *UpdateNodeRequest,
) error {

	// 1. Validate request
	if req == nil {
		return &errors.AnedyaError{
			Message: "update node request cannot be nil",
			Err:     errors.ErrUpdateNodeRequestNil,
		}
	}

	// Validate mandatory NodeID
	if req.NodeID == "" {
		return &errors.AnedyaError{
			Message: "nodeID is required for update",
			Err:     errors.ErrUpdateNodeIDRequired,
		}
	}

	// Ensure at least one update operation is provided
	if len(req.Updates) == 0 {
		return &errors.AnedyaError{
			Message: "at least one update operation must be provided",
			Err:     errors.ErrUpdateNodeEmptyUpdates,
		}
	}

	// Validate each update operation
	for i, u := range req.Updates {

		if u.Type == "" {
			return &errors.AnedyaError{
				Message: fmt.Sprintf("update[%d].type is required", i),
				Err:     fmt.Errorf("update type missing"),
			}
		}

		if u.Type == UpdateTag && u.Tag == nil {
			return &errors.AnedyaError{
				Message: fmt.Sprintf("update[%d].tag is required for tag update", i),
				Err:     fmt.Errorf("tag object missing"),
			}
		}

		if u.Type != UpdateTag && u.Value == "" {
			return &errors.AnedyaError{
				Message: fmt.Sprintf("update[%d].value is required", i),
				Err:     fmt.Errorf("value missing"),
			}
		}
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode update node request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/update", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build update node request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute update node request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read update node response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp UpdateNodeResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode update node response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle API-level errors
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
