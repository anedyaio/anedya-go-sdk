// Package commands provides APIs to issue, track, and manage
// commands sent to nodes within the Anedya platform.
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

// SendCommandRequest represents the request payload used
// to send a command to a node.
type SendCommandRequest struct {
	// NodeId is the unique identifier of the target node.
	NodeId string `json:"nodeId"`

	// Command is the command identifier or name
	// understood by the target device.
	Command string `json:"command"`

	// Data is the command payload sent to the device.
	Data string `json:"data"`

	// Type specifies the data encoding type.
	// Allowed values are "string" or "binary".
	Type CommandDataType `json:"type"`

	// Expiry specifies the command expiration time
	// in seconds from the time of issuance.
	//
	// If omitted, the default expiry is 7 days (604800 seconds).
	Expiry int `json:"expiry,omitempty"`
}

// sendCommandAPIResponse represents the raw response
// returned by the Send Command API.
//
// This type is internal to the SDK and should not be
// exposed to SDK consumers.
type sendCommandAPIResponse struct {
	common.BaseResponse

	// CommandId is the unique identifier assigned
	// to the issued command.
	CommandId string `json:"commandId"`
}

// SendCommandResult represents the processed result
// returned by the SDK after issuing a command.
type SendCommandResult struct {
	// CommandId is the unique identifier
	// of the issued command.
	CommandId string
}

// SendCommand sends a command to a node on the Anedya platform.
//
// If the target node is online and connected via MQTT,
// the command is delivered in real time. Otherwise,
// the command is queued and delivered when the node
// comes back online.
//
// Commands can be queued for a maximum of 7 days,
// after which they are automatically discarded.
// Devices may execute commands sequentially or
// in parallel based on application logic.
//
// Steps performed by this method:
//  1. Validate the request payload and required fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Send Command API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     including cancellation and deadlines.
//   - req: Pointer to SendCommandRequest containing
//     command details.
//
// Returns:
//   - (*SendCommandResult, nil) if the command is sent successfully.
//   - (nil, error) for validation failures, network errors,
//     or API-level errors.
func (cm *CommandManagement) SendCommand(
	ctx context.Context,
	req *SendCommandRequest,
) (*SendCommandResult, error) {

	// validate request
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "send command request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// node ID must be provided
	if req.NodeId == "" {
		return nil, &errors.AnedyaError{
			Message: "nodeId is required",
			Err:     errors.ErrInvalidNodeID,
		}
	}

	// command name must be provided
	if req.Command == "" {
		return nil, &errors.AnedyaError{
			Message: "command is required",
			Err:     errors.ErrCommandRequired,
		}
	}

	// command payload must be provided
	if req.Data == "" {
		return nil, &errors.AnedyaError{
			Message: "data is required",
			Err:     errors.ErrDataRequired,
		}
	}

	// validate command data type
	if req.Type != "string" && req.Type != "binary" {
		return nil, &errors.AnedyaError{
			Message: "type must be either 'string' or 'binary'",
			Err:     errors.ErrInvalidCommandType,
		}
	}

	// validate expiry value
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

	// create HTTP request with context
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

	// execute HTTP request
	resp, err := cm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute send command request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp sendCommandAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode send command response",
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
	return &SendCommandResult{
		CommandId: apiResp.CommandId,
	}, nil
}
