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

// AuthorizeDeviceRequest represents the payload sent to the Authorize Device API.
// It contains identifiers for the node and the device that is being authorized.
type AuthorizeDeviceRequest struct {
	// NodeID is the unique identifier of the node
	// to which the device will be authorized.
	NodeID string `json:"nodeid"`

	// DeviceID is the unique identifier of the device
	// being authorized to connect to the node.
	DeviceID string `json:"deviceid"`
}

// AuthorizeDeviceResponse represents the response returned by
// the Authorize Device API.
//
// It embeds common.BaseResponse to provide standard
// success, error message, and reason code fields.
type AuthorizeDeviceResponse struct {
	common.BaseResponse
}

// AuthorizeDevice authorizes a device to connect to a node in the Anedya platform.
//
// This method performs the following operations:
//  1. Validates the request payload and mandatory fields (NodeID and DeviceID).
//  2. Marshals the request payload into JSON.
//  3. Constructs an HTTP POST request to the Authorize Device API endpoint.
//  4. Executes the HTTP request using the NodeManagement's HTTP client.
//  5. Decodes the API response into AuthorizeDeviceResponse.
//  6. Checks API response status and maps API errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context for controlling request cancellation and timeout.
//   - req: Pointer to AuthorizeDeviceRequest containing node and device identifiers.
//
// Returns:
//   - error: Returns nil if the authorization succeeds, otherwise returns a sentinel
//     error or *errors.AnedyaError if validation, network, or API errors occur.
func (nm *NodeManagement) AuthorizeDevice(
	ctx context.Context,
	req *AuthorizeDeviceRequest,
) error {

	// 1. Validate request
	if req == nil {
		return &errors.AnedyaError{
			Message: "authorize device request cannot be nil",
			Err:     errors.ErrAuthorizeDeviceRequestNil,
		}
	}

	// Validate NodeID
	if req.NodeID == "" {
		return &errors.AnedyaError{
			Message: "nodeId is required to authorize device",
			Err:     errors.ErrAuthorizeDeviceNodeIDRequired,
		}
	}

	// Validate DeviceID
	if req.DeviceID == "" {
		return &errors.AnedyaError{
			Message: "deviceId is required to authorize device",
			Err:     errors.ErrAuthorizeDeviceDeviceIDRequired,
		}
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode authorize device request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/node/authorize", nm.baseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build authorize device request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute authorize device request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read authorize device response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp AuthorizeDeviceResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode authorize device response",
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
