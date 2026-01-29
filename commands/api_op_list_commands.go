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

// ListCommandsRequest represents the request payload used
// to retrieve a list of issued commands.
//
// It supports optional filtering and pagination controls.
type ListCommandsRequest struct {
	// Filter specifies optional criteria for filtering commands,
	// such as node ID, command status, or time range.
	Filter *CommandFilter `json:"filter,omitempty"`

	// Limit specifies the maximum number of commands to return
	// in a single API call (maximum allowed value: 100).
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of commands to skip
	// before starting to return results.
	Offset int `json:"offset,omitempty"`
}

// listCommandsAPIResponse represents the raw response
// returned by the List Commands API.
//
// This type is internal to the SDK and should not be
// exposed to SDK consumers.
type listCommandsAPIResponse struct {
	common.BaseResponse

	// Count is the number of commands returned in this response.
	Count int `json:"count"`

	// TotalCount is the total number of commands
	// matching the applied filter.
	TotalCount int `json:"totalCount"`

	// Data contains the list of command metadata.
	Data []CommandInfo `json:"data"`

	// Next is the pagination token used to fetch
	// the next page of results, if available.
	Next string `json:"next,omitempty"`
}

// ListCommandsResult represents the processed result
// returned by the SDK after listing commands.
type ListCommandsResult struct {
	// Count is the number of commands returned in the current response.
	Count int

	// TotalCount is the total number of commands
	// matching the filter criteria.
	TotalCount int

	// Commands contains the list of retrieved commands.
	Commands []CommandInfo

	// Next is the pagination token for retrieving
	// the next page of results.
	Next string
}

// ListCommands retrieves a list of commands issued to nodes
// with optional filtering and pagination support.
//
// This API returns all commands, including expired ones.
// A maximum of 100 commands can be retrieved in a single request.
//
// Steps performed by this method:
//  1. Validate request parameters.
//  2. Marshal the request payload into JSON.
//  3. Build and send a POST request to the List Commands API.
//  4. Decode the API response.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     including cancellation and timeouts.
//   - req: Pointer to ListCommandsRequest containing
//     filter and pagination options.
//
// Returns:
//   - (*ListCommandsResult, nil) if commands are listed successfully.
//   - (nil, error) for validation failures, network errors,
//     or API-level errors.
func (cm *CommandManagement) ListCommands(
	ctx context.Context,
	req *ListCommandsRequest,
) (*ListCommandsResult, error) {

	// use empty request if nil
	if req == nil {
		req = &ListCommandsRequest{}
	}

	// validate limit
	if req.Limit < 0 || req.Limit > 100 {
		return nil, &errors.AnedyaError{
			Message: "limit must be between 0 and 100",
			Err:     errors.ErrInvalidLimit,
		}
	}

	// validate offset
	if req.Offset < 0 {
		return nil, &errors.AnedyaError{
			Message: "offset must be non-negative",
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

	// execute HTTP request
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
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode list commands response",
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
	return &ListCommandsResult{
		Count:      apiResp.Count,
		TotalCount: apiResp.TotalCount,
		Commands:   apiResp.Data,
		Next:       apiResp.Next,
	}, nil
}
