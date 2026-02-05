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

// CommandFilter represents optional filtering criteria
// used when listing commands from the platform.
type CommandFilter struct {
	// NodeID filters commands issued to a specific node.
	NodeID string `json:"nodeId,omitempty"`

	// Status filters commands by their execution status.
	// Multiple statuses can be provided.
	Status []string `json:"status,omitempty"`
}

// ListCommandsRequest represents the payload used to
// fetch a paginated list of commands from the platform.
type ListCommandsRequest struct {
	// Filter contains optional filtering rules.
	Filter *CommandFilter `json:"filter,omitempty"`

	// Limit specifies the maximum number of commands
	// to return in a single request. Allowed range: 0–100.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the starting index for pagination.
	// Must be greater than or equal to 0.
	Offset int `json:"offset,omitempty"`
}

// listCommandsAPIResponse represents the raw response
// returned by the List Commands API.
//
// This structure is internal to the SDK and is later
// transformed into ListCommandsResult.
type listCommandsAPIResponse struct {
	common.BaseResponse

	// Count indicates the number of commands returned in this page.
	Count int `json:"count"`

	// TotalCount represents the total number of commands available.
	TotalCount int `json:"totalCount"`

	// Data contains the list of command metadata.
	Data []CommandInfo `json:"data"`

	// Next represents the pagination cursor or offset
	// for fetching the next set of results.
	Next int64 `json:"next,omitempty"`
}

// ListCommandsResult represents the processed and
// SDK-friendly result returned to the caller.
type ListCommandsResult struct {
	// Count indicates the number of commands in the current result set.
	Count int

	// TotalCount represents the total number of commands available.
	TotalCount int

	// Commands contains the list of command metadata.
	Commands []CommandInfo

	// Next is the cursor/offset value to be used for the next page request.
	Next int64
}

// ListCommands retrieves a paginated list of commands
// issued within the Anedya platform.
//
// This method supports optional filtering by node ID
// and command status, along with pagination controls.
//
// Steps performed by this method:
//  1. Validate the request payload and pagination constraints.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the List Commands API.
//  4. Read and decode the API response.
//  5. Map API-level errors into structured SDK errors.
//  6. Transform raw API response into ListCommandsResult.
//
// Parameters:
//   - ctx: Context used to control cancellation, timeouts, and deadlines.
//   - req: Pointer to ListCommandsRequest containing filters and pagination options.
//
// Returns:
//   - (*ListCommandsResult, nil) if commands are retrieved successfully.
//   - (nil, error) for validation, network, or decoding failures.
//   - (nil, error) mapped from API-level errors using SDK error mapping.
func (cm *CommandManagement) ListCommands(ctx context.Context, req *ListCommandsRequest) (*ListCommandsResult, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "list commands request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// validate limit range
	if req.Limit < 0 || req.Limit > 100 {
		return nil, &errors.AnedyaError{
			Message: "limit must be between 0 and 100",
			Err:     errors.ErrInvalidLimit,
		}
	}

	// validate offset range
	if req.Offset < 0 {
		return nil, &errors.AnedyaError{
			Message: "offset must be >= 0",
			Err:     errors.ErrInvalidOffset,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/commands/list", cm.baseURL)

	// convert request to JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode list commands request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// create HTTP request with context
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build list commands request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// send HTTP request
	resp, err := cm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute list commands request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read list commands response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// decode API response
	var apiResp listCommandsAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode list commands response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &ListCommandsResult{
		Count:      apiResp.Count,
		TotalCount: apiResp.TotalCount,
		Commands:   apiResp.Data,
		Next:       apiResp.Next,
	}, nil
}
