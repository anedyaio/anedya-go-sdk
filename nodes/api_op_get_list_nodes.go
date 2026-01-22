package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetNodeListRequest represents the payload sent to the Get Node List API.
// It contains pagination and sorting parameters.
type GetNodeListRequest struct {
	// Limit specifies the maximum number of nodes to return in a single request.
	// The value must be between 1 and 1000.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of nodes to skip before starting to return results.
	Offset int `json:"offset,omitempty"`

	// Order specifies the sorting order of the node list.
	// Supported values are:
	//   - "asc"  : ascending order
	//   - "desc" : descending order
	Order string `json:"order,omitempty"`
}

// GetNodeListResponse represents the response returned by the Get Node List API.
type GetNodeListResponse struct {
	// Success indicates whether the request was successful.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`

	// CurrentCount indicates the number of nodes returned in the current response.
	CurrentCount int `json:"currentCount"`

	// TotalCount indicates the total number of nodes available on the platform.
	TotalCount int `json:"totalCount"`

	// Nodes contains the list of node identifiers returned by the API.
	Nodes []string `json:"nodes"`

	// Offset indicates the offset value used for the current response.
	Offset int `json:"offset"`
}

// GetNodeList retrieves a paginated list of nodes from the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory fields (Limit and Order).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Get Node List API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into GetNodeListResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to GetNodeListRequest containing pagination and sorting info.
//
// Returns:
//   - *GetNodeListResponse: Contains node identifiers and pagination metadata on success.
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API errors occur.
func (nm *NodeManagement) GetNodeList(
	ctx context.Context,
	req *GetNodeListRequest,
) (*GetNodeListResponse, error) {

	// Validate request object
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get node list request cannot be nil",
			Err:     errors.ErrNodeListRequestNil,
		}
	}

	// Validate Limit
	if req.Limit <= 0 || req.Limit > 1000 {
		return nil, &errors.AnedyaError{
			Message: "limit must be between 1 and 1000",
			Err:     errors.ErrNodeListInvalidLimit,
		}
	}

	// Validate Order
	if req.Order != "asc" && req.Order != "desc" {
		return nil, &errors.AnedyaError{
			Message: "order must be either 'asc' or 'desc'",
			Err:     errors.ErrNodeListInvalidOrder,
		}
	}

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetNodeList request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/list", nm.baseURL)

	// Build HTTP POST request with context
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build GetNodeList request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetNodeList request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode API response
	var apiResp GetNodeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetNodeList response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// API-level error handling
	if !apiResp.Success {
		sdkErr := errors.GetError(apiResp.ReasonCode, apiResp.Error)
		// Return any other API errors
		return nil, sdkErr
	}

	// Success
	return &apiResp, nil
}
