// Package nodes provides APIs to manage nodes and retrieve
// node-related information from the Anedya platform.
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

// GetNodeDetailsRequest represents the payload used to
// fetch detailed information for one or more nodes.
type GetNodeDetailsRequest struct {
	// Nodes is the list of node IDs for which details
	// are to be retrieved.
	Nodes []string `json:"nodes"`
}

// getNodeDetailsAPIResponse represents the raw response
// returned by the Node Details API.
//
// This structure is internal and mapped directly to
// the SDK result before being returned to the caller.
type getNodeDetailsAPIResponse struct {
	common.BaseResponse

	// Data maps node IDs to their corresponding
	// node details.
	Data map[string]Node `json:"data,omitempty"`
}

// GetNodeDetails retrieves detailed information for one
// or more nodes from the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Node Details API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and deadlines.
//   - req: Pointer to GetNodeDetailsRequest containing
//     the list of node IDs.
//
// Returns:
//   - (map[string]Node, nil) on successful execution.
//   - (nil, error) for validation, network, or API errors.
func (nm *NodeManagement) GetNodeDetails(
	ctx context.Context,
	req *GetNodeDetailsRequest,
) (map[string]Node, error) {

	// 1. Validate request
	if req == nil || len(req.Nodes) == 0 {
		return nil, &errors.AnedyaError{
			Message: "node list cannot be empty",
			Err:     errors.ErrNodeDetailsRequestNil,
		}
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode get node details request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/details", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build get node details request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute get node details request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read get node details response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp getNodeDetailsAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode get node details response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Success
	return apiResp.Data, nil
}
