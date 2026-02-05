// Package commands provides APIs to issue commands to nodes
// and retrieve their execution status within the Anedya platform.
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

// GetCommandStatusRequest represents the payload used to fetch
// the current status and details of a previously issued command.
type GetCommandStatusRequest struct {
	// CommandId is the unique identifier of the command
	// whose status is being requested.
	CommandId string `json:"commandId"`
}

// getCommandStatusAPIResponse represents the raw response
// returned by the Get Command Details API.
//
// This structure is internal to the SDK and is later
// transformed into a user-friendly result type.
type getCommandStatusAPIResponse struct {
	common.BaseResponse

	CommandId   string `json:"commandId"`
	Command     string `json:"command"`
	Status      string `json:"status"`
	UpdatedOn   int64  `json:"updatedOn"`
	AckData     string `json:"ackdata"`
	AckDataType string `json:"ackdatatype"`
	Expired     bool   `json:"expired"`
	Expiry      int64  `json:"expiry"`
	IssuedAt    int64  `json:"issuedAt"`
	Data        string `json:"data"`
	DataType    string `json:"datatype"`
}

// GetCommandStatusResult represents the processed and
// SDK-friendly result returned after fetching command status.
type GetCommandStatusResult struct {
	// Command contains all command metadata,
	// acknowledgement data, and payload details.
	Command CommandDetails
}

// GetCommandStatus retrieves the execution status and metadata
// of a command previously issued to a node in the Anedya platform.
//
// Steps performed by this method:
//  1. Validate the request payload and mandatory fields.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Command Details API.
//  4. Read and decode the API response.
//  5. Map API-level errors into structured SDK errors.
//  6. Transform raw API response into GetCommandStatusResult.
//
// Parameters:
//   - ctx: Context used to control cancellation, timeouts, and deadlines.
//   - req: Pointer to GetCommandStatusRequest containing the command ID.
//
// Returns:
//   - (*GetCommandStatusResult, nil) if the command status is retrieved successfully.
//   - (nil, error) for validation, network, or decoding failures.
//   - (nil, error) mapped from API-level errors using SDK error mapping.
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
	if req.CommandId == "" {
		return nil, &errors.AnedyaError{
			Message: "command id is required",
			Err:     errors.ErrInvalidCommandID,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/commands/getDetails", cm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode get command status request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// create HTTP request with context
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read get command status response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// decode API response
	var apiResp getCommandStatusAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
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
		Command: CommandDetails{
			CommandInfo: CommandInfo{
				CommandId: apiResp.CommandId,
				Command:   apiResp.Command,
				Status:    CommandStatus(apiResp.Status),
				UpdatedOn: apiResp.UpdatedOn,
				Expired:   apiResp.Expired,
				Expiry:    apiResp.Expiry,
				IssuedAt:  apiResp.IssuedAt,
			},
			AckData:     apiResp.AckData,
			AckDataType: CommandDataType(apiResp.AckDataType),
			Data:        apiResp.Data,
			DataType:    CommandDataType(apiResp.DataType),
		},
	}, nil
}
