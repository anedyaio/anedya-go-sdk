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

// ClearChildNodesRequest represents the payload sent to the Clear Child Nodes API.
// It contains the parent node identifier whose child nodes are to be removed.
type ClearChildNodesRequest struct {
	// ParentId is the unique identifier of the parent node
	// whose child nodes will be removed.
	ParentId string `json:"parentId"`
}

// ClearChildNodesResponse represents the response returned by
// the Clear Child Nodes API.
//
// It embeds common.BaseResponse to provide standard
// success, error message, and reason code fields.
type ClearChildNodesResponse struct {
	common.BaseResponse
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

	// 1. Validate request
	if req == nil {
		return &errors.AnedyaError{
			Message: "clear child nodes request cannot be nil",
			Err:     errors.ErrClearChildNodesRequestNil,
		}
	}

	if req.ParentId == "" {
		return &errors.AnedyaError{
			Message: "parentId is required to clear child nodes",
			Err:     errors.ErrClearChildNodesParentIDRequired,
		}
	}

	// 2. Encode request
	requestBody, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode clear child nodes request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/child/clear", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build clear child nodes request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute clear child nodes request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read clear child nodes response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp ClearChildNodesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode clear child nodes response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle API-level errors
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Success
	return nil
}
