// Package commands provides APIs to issue commands to nodes
// and retrieve command execution status from the Anedya platform.
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
// retrieve the status of a previously issued command.
type GetCommandStatusRequest struct {
	// Id is the unique identifier of the command
	// whose status is being requested.
	Id string `json:"id"`
}

// getCommandStatusAPIResponse represents the raw response
// returned by the Command Status API.
//
// This structure is internal and mapped to
// GetCommandStatusResult before being returned to the caller.
type getCommandStatusAPIResponse struct {
	common.BaseResponse
	CommandDetails
}

// GetCommandStatusResult represents the processed result
// returned to SDK consumers after fetching command status.
type GetCommandStatusResult struct {
	// Command contains the detailed information
	// associated with the requested command.
	Command CommandDetails
}

// GetCommandStatus retrieves the current execution status
// and detailed metadata of a command.
//
// This method returns comprehensive information including:
//   - Current command status (pending, received, processing,
//     success, failure, invalidated)
//   - Command issue, update, and expiry timestamps
//   - Command payload and acknowledgment data
//   - Expiration and validity status
//
// Steps performed by this method:
//  1. Validate the request payload.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Command Status API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and deadlines.
//   - req: Pointer to GetCommandStatusRequest containing
//     the command identifier.
//
// Returns:
//   - (*GetCommandStatusResult, nil) if the command status
//     is retrieved successfully.
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

	// validate command ID
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

	// execute HTTP request
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
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode get command status response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle HTTP-level or API-level errors
	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices ||
		!apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &GetCommandStatusResult{
		Command: apiResp.CommandDetails,
	}, nil
}
