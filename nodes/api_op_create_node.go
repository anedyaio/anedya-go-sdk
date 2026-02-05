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

// CreateNodeRequest represents the payload sent to the Create Node API.
// It contains the mandatory node name and optional metadata.
type CreateNodeRequest struct {
	// NodeName is the human-readable name of the node.
	// This field is mandatory.
	NodeName string `json:"node_name"`

	// NodeDesc provides an optional description of the node.
	NodeDesc string `json:"node_desc,omitempty"`

	// Tags contains optional metadata tags associated with the node.
	Tags []Tag `json:"tags,omitempty"`

	// PreauthId optionally associates the node with a pre-authorized identifier.
	PreauthId string `json:"preauth_id,omitempty"`
}

// CreateNodeResponse represents the response returned by the Create Node API.
type CreateNodeResponse struct {
	common.BaseResponse

	// NodeId is the unique identifier assigned to the newly created node.
	NodeId string `json:"nodeId,omitempty"`
}

// CreateNode creates a new node in the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory fields (NodeName).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Create Node API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into CreateNodeResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to CreateNodeRequest containing node information.
//
// Returns:
//   - *Node: A Node object representing the newly created node on success.
//   - error: Returns nil on success, otherwise a sentinel error or *errors.AnedyaError
//     if validation, network, or API errors occur.
func (nm *NodeManagement) CreateNode(
	ctx context.Context,
	req *CreateNodeRequest,
) (*Node, error) {

	// 1. Validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "create node request cannot be nil",
			Err:     errors.ErrInputRequired,
		}
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode create node request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/create", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build create node request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute create node request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read create node response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp CreateNodeResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode create node response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Success: return Node
	return &Node{
		NodeId:          apiResp.NodeId,
		NodeName:        req.NodeName,
		NodeDescription: req.NodeDesc,
		Tags:            req.Tags,
		nodeManagement:  nm,
	}, nil
}
