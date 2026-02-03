// Package commands provides APIs to manage and query commands
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

// ListCommandsRequest represents the payload used to
// fetch a paginated list of commands with optional filters.
type ListCommandsRequest struct {
	// Filter contains optional filtering criteria
	// such as node, status, or time range.
	Filter *CommandFilter `json:"filter,omitempty"`

	// Limit defines the maximum number of commands
	// to return in a single request.
	// Valid range: 0–100.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the starting index
	// for pagination.
	Offset int `json:"offset,omitempty"`
}

// listCommandsAPIResponse represents the raw API response
// returned by the list commands endpoint.
type listCommandsAPIResponse struct {
	common.BaseResponse

	// Count is the number of commands returned
	// in the current response.
	Count int `json:"count"`

	// TotalCount is the total number of commands
	// available that match the filter.
	TotalCount int `json:"totalCount"`

	// Data contains the list of command metadata.
	Data []CommandInfo `json:"data"`

	// Next is the pagination cursor or token
	// used to fetch the next set of results.
	Next string `json:"next,omitempty"`
}

// ListCommandsResult represents the processed result
// returned to SDK consumers.
type ListCommandsResult struct {
	// Count is the number of commands returned
	// in this result set.
	Count int

	// TotalCount is the total number of matching commands.
	TotalCount int

	// Commands contains command details.
	Commands []CommandInfo

	// Next is the pagination cursor/token
	// for the next request.
	Next string
}

// ListCommands retrieves a paginated list of commands
// from the Anedya platform with optional filtering.
//
// Steps performed by this method:
//  1. Validate the request payload and pagination parameters.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the list commands API.
//  4. Decode the API response.
//  5. Convert API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and timeouts.
//   - req: Pointer to ListCommandsRequest containing
//     filter and pagination options.
//
// Returns:
//   - (*ListCommandsResult, nil) if successful.
//   - (nil, error) for validation, network, or API errors.
func (cm *CommandManagement) ListCommands(
	ctx context.Context,
	req *ListCommandsRequest,
) (*ListCommandsResult, error) {

	// request must not be nil
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

	// validate offset
	if req.Offset < 0 {
		return nil, &errors.AnedyaError{
			Message: "offset must be >= 0",
			Err:     errors.ErrInvalidOffset,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/commands/list", cm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode list commands request",
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

	// decode API response
	var apiResp listCommandsAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
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
