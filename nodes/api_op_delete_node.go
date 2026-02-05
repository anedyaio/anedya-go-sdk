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

// DeleteNodeRequest represents the payload sent to the Delete Node API.
// It contains the unique identifier of the node that needs to be deleted.
type DeleteNodeRequest struct {
	// NodeID is the unique identifier of the node to be deleted.
	NodeID string `json:"nodeid"`
}

// DeleteNodeResponse represents the response returned by the Delete Node API.
type DeleteNodeResponse struct {
	common.BaseResponse
}

// DeleteNode deletes a node from the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory field (NodeID).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Delete Node API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into DeleteNodeResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to DeleteNodeRequest containing the node identifier.
//
// Returns:
//   - error: Returns nil if the operation succeeds, otherwise a sentinel
//     error or *errors.AnedyaError if validation, network, or API errors occur.
func (nm *NodeManagement) DeleteNode(
	ctx context.Context,
	req *DeleteNodeRequest,
) error {

	// 1. Validate request
	if req == nil {
		return &errors.AnedyaError{
			Message: "delete node request cannot be nil",
			Err:     errors.ErrDeleteNodeRequestNil,
		}
	}

	if req.NodeID == "" {
		return &errors.AnedyaError{
			Message: "node id is required to delete node",
			Err:     errors.ErrDeleteNodeIDRequired,
		}
	}

	// 2. Encode request
	requestBody, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode delete node request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/delete", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build delete node request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute delete node request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read delete node response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp DeleteNodeResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode delete node response",
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
