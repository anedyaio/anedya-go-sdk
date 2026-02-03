// Package commands provides APIs to manage commands,
// including issuing, tracking, and invalidating commands
// within the Anedya platform.
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
// invalidate or cancel a previously issued command.
type InvalidateCommandRequest struct {
	// Id is the unique identifier of the command
	// to be invalidated.
	Id string `json:"id"`
}

// invalidateCommandAPIResponse represents the raw API response
// returned by the invalidate command endpoint.
type invalidateCommandAPIResponse struct {
	common.BaseResponse
}

// InvalidateCommandResult represents the processed result
// returned to SDK consumers after invalidation.
type InvalidateCommandResult struct {
	// Success indicates whether the command
	// was successfully invalidated.
	Success bool
}

// InvalidateCommand invalidates or cancels an existing command
// in the Anedya platform so it will no longer be executed.
//
// Steps performed by this method:
//  1. Validate the request payload and command identifier.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the invalidate command API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and timeouts.
//   - req: Pointer to InvalidateCommandRequest containing
//     the command identifier.
//
// Returns:
//   - (*InvalidateCommandResult, nil) if invalidation succeeds.
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

	// command ID must be provided
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

	// send HTTP request
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
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode invalidate command response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &InvalidateCommandResult{
		Success: true,
	}, nil
}
