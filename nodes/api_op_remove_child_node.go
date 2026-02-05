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

// RemoveChildNodeRequest represents the request payload sent to the Remove Child Node API.
// It contains the parent node ID and the child node ID to be detached.
type RemoveChildNodeRequest struct {
	// ParentId is the unique identifier of the parent node.
	ParentId string `json:"parentId"`

	// ChildNode is the unique identifier of the child node
	// to be removed from the parent.
	ChildNode string `json:"childNode"`
}

// RemoveChildNodeResponse represents the response returned by the Remove Child Node API.
type RemoveChildNodeResponse struct {
	common.BaseResponse
}

// RemoveChildNode detaches a child node from its parent node in the Anedya platform.
//
// The request is provided using a *RemoveChildNodeRequest structure,
// which specifies both the parent node ID and the child node ID to remove.
//
// The method performs the following operations:
//  1. Validates the request payload and mandatory fields.
//  2. Marshals the request into JSON.
//  3. Constructs an HTTP POST request to the Remove Child Node API endpoint.
//  4. Executes the request using NodeManagement's HTTP client.
//  5. Decodes the API response.
//  6. Maps API or network errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - req: Pointer to RemoveChildNodeRequest containing parent and child node IDs.
//
// Returns:
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API-level errors occur.
func (nm *NodeManagement) RemoveChildNode(
	ctx context.Context,
	req *RemoveChildNodeRequest,
) error {

	// 1. Validate request
	if req == nil {
		return &errors.AnedyaError{
			Message: "remove child node request cannot be nil",
			Err:     errors.ErrRemoveChildNodeRequestNil,
		}
	}
	if req.ParentId == "" {
		return &errors.AnedyaError{
			Message: "parent id is required",
			Err:     errors.ErrRemoveChildNodeParentIDRequired,
		}
	}
	if req.ChildNode == "" {
		return &errors.AnedyaError{
			Message: "child node id is required",
			Err:     errors.ErrRemoveChildNodeChildIDRequired,
		}
	}

	// 2. Encode request
	requestBody, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode remove child node request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/child/remove", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build remove child node request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute remove child node request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read remove child node response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp RemoveChildNodeResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode remove child node response",
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
