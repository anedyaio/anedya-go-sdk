// Package commands provides APIs for issuing, tracking,
// and managing device commands within the Anedya platform.
package commands

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

// InvalidateCommandRequest represents the payload used to
// invalidate or cancel a previously issued command.
type InvalidateCommandRequest struct {
	// CommandId is the unique identifier of the command
	// that needs to be invalidated.
	CommandId string `json:"commandId"`
}

// invalidateCommandAPIResponse represents the raw API response
// received from the backend before SDK transformation.
type invalidateCommandAPIResponse struct {
	common.BaseResponse
}

// InvalidateCommandResult represents the processed result
// returned to SDK consumers after successful invalidation.
type InvalidateCommandResult struct {
	// Success indicates whether the command
	// was successfully invalidated.
	Success bool
}

// InvalidateCommand cancels or invalidates a previously issued command.
//
// Workflow:
//  1. Validate request object and command identifier.
//  2. Marshal request into JSON.
//  3. Send POST request to invalidate endpoint.
//  4. Decode API response.
//  5. Return structured SDK result.
//
// Parameters:
//   - ctx: Context used for cancellation and timeout control.
//   - req: Pointer to InvalidateCommandRequest containing command ID.
//
// Returns:
//   - (*InvalidateCommandResult, nil) if invalidation succeeds.
//   - (nil, error) if validation, network, decoding,
//     or API-level errors occur.
func (cm *CommandManagement) InvalidateCommand(
	ctx context.Context,
	req *InvalidateCommandRequest,
) (*InvalidateCommandResult, error) {

	// validate request object
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "invalidate command request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// validate command ID
	if req.CommandId == "" {
		return nil, &errors.AnedyaError{
			Message: "command id is required",
			Err:     errors.ErrInvalidCommandID,
		}
	}

	// construct API endpoint URL
	url := fmt.Sprintf("%s/v1/commands/invalidate", cm.baseURL)

	// encode request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode invalidate command request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// build HTTP request
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

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read invalidate command response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// decode JSON response
	var apiResp invalidateCommandAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode invalidate command response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle API-level failure
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// return structured result
	return &InvalidateCommandResult{
		Success: true,
	}, nil
}
