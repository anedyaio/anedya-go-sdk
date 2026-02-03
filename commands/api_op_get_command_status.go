// Package commands provides APIs to manage device or node
// commands and retrieve their execution status within
// the Anedya platform.
package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// GetCommandStatusRequest represents the payload used to
// request the status of a previously issued command.
type GetCommandStatusRequest struct {
	// Id is the unique identifier of the command
	// whose status is being requested.
	Id string `json:"id"`
}

// getCommandStatusAPIResponse represents the raw API response
// returned by the command status endpoint.
type getCommandStatusAPIResponse struct {
	common.BaseResponse

	// CommandDetails contains detailed metadata and
	// execution status of the command.
	CommandDetails
}

// GetCommandStatusResult represents the processed result
// returned to SDK consumers.
type GetCommandStatusResult struct {
	// Command contains the command metadata and
	// its current execution status.
	Command CommandDetails
}

// GetCommandStatus retrieves the execution status and details
// of a command from the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload and command identifier.
//  2. Marshal the request into JSON.
//  3. Build and send a POST request to the command status API.
//  4. Decode the API response into an internal response struct.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and timeouts.
//   - req: Pointer to GetCommandStatusRequest containing
//     the command identifier.
//
// Returns:
//   - (*GetCommandStatusResult, nil) on successful retrieval.
//   - (nil, error) for validation, network, or API errors.
func (cm *CommandManagement) GetCommandStatus(
	ctx context.Context,
	req *GetCommandStatusRequest,
) (*GetCommandStatusResult, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "get command status request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// command ID must be provided
	if req.Id == "" {
		return nil, &errors.AnedyaError{
			Message: "command id is required",
			Err:     errors.ErrInvalidCommandID,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/commands/status", cm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode get command status request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// create HTTP request with context
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build get command status request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// send HTTP request
	resp, err := cm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute get command status request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp getCommandStatusAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode get command status response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &GetCommandStatusResult{
		Command: apiResp.CommandDetails,
	}, nil
}
