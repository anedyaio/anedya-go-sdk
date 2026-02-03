// Package nodes provides APIs to manage nodes and retrieve
// node-related information from the Anedya platform.
package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetNodeListRequest represents the payload used to
// retrieve a paginated list of nodes.
//
// It supports pagination and ordering of node IDs.
type GetNodeListRequest struct {
	// Limit specifies the maximum number of nodes
	// to return in a single request.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the starting position
	// for pagination.
	Offset int `json:"offset,omitempty"`

	// Order specifies the sorting order of nodes.
	// Allowed values are "asc" or "desc".
	Order string `json:"order,omitempty"`
}

// getNodeListAPIResponse represents the raw response
// returned by the Node List API.
//
// This structure is internal and mapped to
// GetNodeListResult before being returned to the caller.
type getNodeListAPIResponse struct {
	common.BaseResponse

	// CurrentCount is the number of nodes returned
	// in the current response.
	CurrentCount int `json:"currentCount"`

	// TotalCount is the total number of nodes
	// available in the system.
	TotalCount int `json:"totalCount"`

	// Nodes contains the list of node IDs.
	Nodes []string `json:"nodes"`

	// Offset indicates the pagination offset
	// used for this response.
	Offset int `json:"offset"`
}

// GetNodeListResult represents the processed result
// returned to SDK consumers after fetching the node list.
type GetNodeListResult struct {
	// CurrentCount is the number of nodes
	// returned in the current response.
	CurrentCount int

	// TotalCount is the total number of nodes
	// available across all pages.
	TotalCount int

	// Offset indicates the pagination offset
	// used to retrieve this result.
	Offset int

	// Nodes contains the list of node IDs.
	Nodes []string
}

// GetNodeList retrieves a paginated list of node IDs
// from the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload and parameters.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Node List API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//  6. Convert the API response into GetNodeListResult.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and deadlines.
//   - req: Pointer to GetNodeListRequest containing
//     pagination and ordering options.
//
// Returns:
//   - (*GetNodeListResult, nil) on successful execution.
//   - (nil, error) for validation, network, or API errors.
func (nm *NodeManagement) GetNodeList(
	ctx context.Context,
	req *GetNodeListRequest,
) (*GetNodeListResult, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get node list request cannot be nil",
			Err:     errors.ErrNodeListRequestNil,
		}
	}

	// validate limit range
	if req.Limit <= 0 || req.Limit > 1000 {
		return nil, &errors.AnedyaError{
			Message: "limit must be between 1 and 1000",
			Err:     errors.ErrNodeListInvalidLimit,
		}
	}

	// validate order field
	if req.Order != "" && req.Order != "asc" && req.Order != "desc" {
		return nil, &errors.AnedyaError{
			Message: "order must be either 'asc' or 'desc'",
			Err:     errors.ErrNodeListInvalidOrder,
		}
	}

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetNodeList request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/node/list", nm.baseURL)

	// create HTTP request with context
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build GetNodeList request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// send HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetNodeList request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp getNodeListAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetNodeList response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// API-level error handling
	if !apiResp.Success {
		sdkErr := errors.GetError(apiResp.ReasonCode, apiResp.Error)
		// Return any other API errors
		return nil, sdkErr
	}

	// Success
	return &GetNodeListResult{
		CurrentCount: apiResp.CurrentCount,
		TotalCount:   apiResp.TotalCount,
		Offset:       apiResp.Offset,
		Nodes:        apiResp.Nodes,
	}, nil
}
