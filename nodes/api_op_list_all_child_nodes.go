// Package nodes provides APIs to manage nodes and
// retrieve node hierarchy information from Anedya.
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

// ListChildNodesRequest represents the payload used to
// retrieve child nodes of a given parent node.
type ListChildNodesRequest struct {
	// ParentId is the unique identifier of the parent node.
	// This field is required.
	ParentId string `json:"parentId"`

	// Limit specifies the maximum number of child nodes
	// to return in a single request.
	//
	// If not provided or out of range, a default value
	// is applied by the SDK.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the starting point for pagination.
	// It is used to fetch the next set of child nodes.
	Offset int `json:"offset,omitempty"`
}

// ChildNode represents a direct child node associated
// with a parent node.
type ChildNode struct {
	// ChildId is the unique identifier of the child node.
	ChildId string `json:"childId"`

	// Alias is the human-readable alias of the child node.
	Alias string `json:"alias"`

	// CreatedAt is the timestamp (in milliseconds since epoch)
	// when the child node was created.
	CreatedAt int64 `json:"createdAt"`
}

// listChildNodesAPIResponse represents the raw response
// returned by the List Child Nodes API.
//
// This structure is internal and converted into
// ListChildNodesResult before being returned to the caller.
type listChildNodesAPIResponse struct {
	common.BaseResponse

	// TotalCount represents the total number of child nodes
	// available for the given parent node.
	TotalCount int `json:"totalCount"`

	// Count represents the number of child nodes returned
	// in the current response.
	Count int `json:"count"`

	// Next represents the offset value to be used
	// for fetching the next page of results.
	Next int `json:"next"`

	// Data contains the list of child nodes.
	Data []ChildNode `json:"data"`
}

// ListChildNodesResult represents the structured result
// returned after successfully fetching child nodes.
type ListChildNodesResult struct {
	// TotalCount is the total number of child nodes
	// associated with the parent node.
	TotalCount int

	// Count is the number of child nodes returned
	// in the current request.
	Count int

	// Next is the offset to be used for the next
	// pagination request.
	Next int

	// Nodes contains the list of retrieved child nodes.
	Nodes []ChildNode
}

// ListChildNodes retrieves the list of direct child nodes
// for a given parent node.
//
// This method supports pagination using limit and offset.
// If limit or offset values are invalid, the SDK applies
// sensible defaults.
//
// Steps performed by this method:
//  1. Validate the request payload.
//  2. Normalize pagination parameters.
//  3. Marshal the request into JSON.
//  4. Build and execute the HTTP request.
//  5. Decode and validate the API response.
//
// Parameters:
//   - ctx: Context used to manage request lifecycle,
//     cancellation, and deadlines.
//   - req: Pointer to ListChildNodesRequest containing
//     parent node ID and pagination options.
//
// Returns:
//   - (*ListChildNodesResult, nil) on success.
//   - (nil, error) for validation, network, or API errors.
func (nm *NodeManagement) ListChildNodes(
	ctx context.Context,
	req *ListChildNodesRequest,
) (*ListChildNodesResult, error) {

	// 1. Validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "list child nodes request cannot be nil",
			Err:     errors.ErrListChildNodesRequestNil,
		}
	}

	// validate parent node ID
	if req.ParentId == "" {
		return nil, &errors.AnedyaError{
			Message: "parent id is required to list child nodes",
			Err:     errors.ErrListChildNodesParentIDRequired,
		}
	}

	// normalize pagination values
	if req.Limit <= 0 || req.Limit > 1000 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode list child nodes request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/child/list", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build list child nodes request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute list child nodes request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read list child nodes response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp listChildNodesAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode list child nodes response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. API-level error handling
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Success
	return &ListChildNodesResult{
		TotalCount: apiResp.TotalCount,
		Count:      apiResp.Count,
		Next:       apiResp.Next,
		Nodes:      apiResp.Data,
	}, nil
}
