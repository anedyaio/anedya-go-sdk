// Package commands provides APIs to send, manage,
// and query device commands in the Anedya platform.
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

// SendCommandRequest represents the payload used to
// send a command to a specific node/device.
type SendCommandRequest struct {
	// NodeId is the unique identifier of the target node/device.
	NodeId string `json:"nodeId"`

	// Command is the name or action to be executed on the node.
	Command string `json:"command"`

	// Data is the command payload.
	// It can be a string or base64-encoded binary depending on Type.
	Data string `json:"data"`

	// Type specifies the format of Data.
	// Supported values: "string", "binary".
	Type CommandDataType `json:"type"`

	// Expiry defines the time (in seconds) after which
	// the command becomes invalid if not executed.
	// Optional.
	Expiry int `json:"expiry,omitempty"`
}

// sendCommandAPIResponse represents the raw API response
// returned by the send command endpoint.
type sendCommandAPIResponse struct {
	common.BaseResponse

	// CommandId is the unique identifier
	// generated for the sent command.
	CommandId string `json:"commandId"`
}

// SendCommandResult represents the processed SDK result
// returned to consumers after sending a command.
type SendCommandResult struct {
	// CommandId uniquely identifies the created command.
	CommandId string
}

// SendCommand sends a command to a specific node/device
// in the Anedya platform.
//
// Steps performed by this method:
//  1. Validate request payload fields.
//  2. Marshal request into JSON.
//  3. Build and execute HTTP POST request.
//  4. Decode API response.
//  5. Convert API-level errors into SDK errors.
//
// Parameters:
//   - ctx: Context controlling cancellation and timeout.
//   - req: Pointer to SendCommandRequest containing command details.
//
// Validation Rules:
//   - NodeId must not be empty.
//   - Command must not be empty.
//   - Data must not be empty.
//   - Type must be "string" or "binary".
//   - Expiry must be >= 0.
//
// Returns:
//   - (*SendCommandResult, nil) on success.
//   - (nil, error) on validation, network, or API errors.
func (cm *CommandManagement) SendCommand(
	ctx context.Context,
	req *SendCommandRequest,
) (*SendCommandResult, error) {

	// 1. Validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "send command request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	if req.NodeId == "" {
		return nil, &errors.AnedyaError{
			Message: "nodeId is required",
			Err:     errors.ErrInvalidNodeID,
		}
	}

	if req.Command == "" {
		return nil, &errors.AnedyaError{
			Message: "command is required",
			Err:     errors.ErrCommandRequired,
		}
	}

	if req.Data == "" {
		return nil, &errors.AnedyaError{
			Message: "data is required",
			Err:     errors.ErrDataRequired,
		}
	}

	if req.Type != "string" && req.Type != "binary" {
		return nil, &errors.AnedyaError{
			Message: "type must be either 'string' or 'binary'",
			Err:     errors.ErrInvalidCommandType,
		}
	}

	if req.Expiry < 0 {
		return nil, &errors.AnedyaError{
			Message: "expiry must be a positive value",
			Err:     errors.ErrInvalidExpiry,
		}
	}

	// 2. Build API URL
	url := fmt.Sprintf("%s/v1/commands/send", cm.baseURL)

	// 3. Marshal request to JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode send command request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 4. Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build send command request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 5. Execute request
	resp, err := cm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute send command request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 6. Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read send command response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 7. Decode response
	var apiResp sendCommandAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode send command response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 8. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 9. Success
	return &SendCommandResult{
		CommandId: apiResp.CommandId,
	}, nil
}
