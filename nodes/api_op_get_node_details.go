package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetNodeDetailsRequest represents the request payload sent to the Get Node Details API.
// It contains a list of node IDs whose details are requested.
type GetNodeDetailsRequest struct {
	// Nodes is a list of node IDs for which details are required.
	// At least one node ID must be provided.
	Nodes []string `json:"nodes"`
}

// GetNodeDetailsResponse represents the response returned by the Get Node Details API.
type GetNodeDetailsResponse struct {
	// Success indicates whether the request was processed successfully.
	Success bool `json:"success"`

	// Error contains a human-readable error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode is a machine-readable error code used for SDK error mapping.
	ReasonCode string `json:"reasonCode,omitempty"`

	// Data contains node details mapped by node ID.
	Data map[string]Node `json:"data,omitempty"`
}

// GetNodeDetails retrieves detailed information for one or more nodes from the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and ensures at least one node ID is provided.
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Get Node Details API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into GetNodeDetailsResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to GetNodeDetailsRequest containing the node IDs.
//
// Returns:
//   - map[string]NodeDetails: Mapping of node IDs to NodeDetails on success.
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API errors occur.
func (nm *NodeManagement) GetNodeDetails(
	ctx context.Context,
	req *GetNodeDetailsRequest,
) (map[string]Node, error) {

	// Validate request object and ensure at least one node ID is provided
	if req == nil || len(req.Nodes) == 0 {
		return nil, errors.ErrNodeListRequestNil
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/details", nm.baseURL)

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode GetNodeDetails request",
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
			Message: "failed to build GetNodeDetails request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute GetNodeDetails request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON into GetNodeDetailsResponse
	var apiResp GetNodeDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode GetNodeDetails response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Check HTTP status code for success
	if resp.StatusCode != http.StatusOK {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Check API-level success flag
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Return the node details map
	return apiResp.Data, nil
}
