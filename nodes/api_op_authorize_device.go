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

	// Validate request object
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

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/authorize", nm.baseURL)

	// Marshal request payload to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode AuthorizeDevice request",
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
		return &errors.AnedyaError{
			Message: "failed to build AuthorizeDevice request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute AuthorizeDevice request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Decode response JSON
	var apiResp AuthorizeDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode AuthorizeDevice response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle HTTP or API level error
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
