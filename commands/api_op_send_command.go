// Package commands provides APIs to send, manage,
// and query device commands in the Anedya platform.
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

	// request must not be nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "send command request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// validate node ID
	if req.NodeId == "" {
		return nil, &errors.AnedyaError{
			Message: "nodeId is required",
			Err:     errors.ErrInvalidNodeID,
		}
	}

	// validate command name
	if req.Command == "" {
		return nil, &errors.AnedyaError{
			Message: "command is required",
			Err:     errors.ErrCommandRequired,
		}
	}

	// validate payload
	if req.Data == "" {
		return nil, &errors.AnedyaError{
			Message: "data is required",
			Err:     errors.ErrDataRequired,
		}
	}

	// validate data type
	if req.Type != "string" && req.Type != "binary" {
		return nil, &errors.AnedyaError{
			Message: "type must be either 'string' or 'binary'",
			Err:     errors.ErrInvalidCommandType,
		}
	}

	// validate expiry
	if req.Expiry < 0 {
		return nil, &errors.AnedyaError{
			Message: "expiry must be a positive value",
			Err:     errors.ErrInvalidExpiry,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/commands/send", cm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode send command request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build send command request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// execute request
	resp, err := cm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute send command request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode response
	var apiResp sendCommandAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode send command response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// API-level error handling
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &SendCommandResult{
		CommandId: apiResp.CommandId,
	}, nil
}
