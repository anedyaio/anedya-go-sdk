// Package commands provides APIs to manage device commands
// within the Anedya platform, including issuing, tracking,
// and invalidating commands.
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

// InvalidateCommandRequest represents the payload used to
// invalidate a previously issued command.
type InvalidateCommandRequest struct {
	// Id is the unique identifier of the command to invalidate.
	Id string `json:"id"`
}

// invalidateCommandAPIResponse represents the raw response
// returned by the Invalidate Command API.
//
// This structure is internal and used only for decoding
// the API response.
type invalidateCommandAPIResponse struct {
	common.BaseResponse
}

// InvalidateCommandResult represents the result returned
// to SDK consumers after attempting to invalidate a command.
type InvalidateCommandResult struct {
	// Success indicates whether the command was successfully invalidated.
	Success bool
}

// InvalidateCommand cancels a pending command issued to a device.
//
// A command can only be invalidated if its status is "pending".
// Once the command has been received or is being processed
// by the device, it can no longer be invalidated.
//
// This method is useful when:
//   - A user cancels an operation before device execution
//   - A newer command supersedes an older queued command
//   - An incorrect command needs to be withdrawn
//
// Steps performed by this method:
//  1. Validate the request payload and command ID.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Invalidate Command API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and deadlines.
//   - req: Pointer to InvalidateCommandRequest containing
//     the command identifier.
//
// Returns:
//   - (*InvalidateCommandResult, nil) if the command is invalidated successfully.
//   - (nil, error) for validation, network, or API errors.
func (cm *CommandManagement) InvalidateCommand(
	ctx context.Context,
	req *InvalidateCommandRequest,
) (*InvalidateCommandResult, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "invalidate command request cannot be nil",
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
	url := fmt.Sprintf("%s/v1/commands/invalidate", cm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode invalidate command request",
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
			Message: "failed to build invalidate command request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// execute HTTP request
	resp, err := cm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute invalidate command request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp invalidateCommandAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode invalidate command response",
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
	return &InvalidateCommandResult{
		Success: true,
	}, nil
}
