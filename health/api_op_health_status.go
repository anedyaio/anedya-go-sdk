// Package health provides APIs to check device connectivity
// and monitor the last heartbeat status of nodes
// on the Anedya platform.
package health

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

const (
	// maxHealthThresholdSec defines the maximum allowed
	// heartbeat threshold duration (7 days in seconds).
	maxHealthThresholdSec = 7 * 24 * 60 * 60
)

// HealthStatusRequest represents the payload used to query
// the health status of one or more nodes.
//
// The health status is determined based on the last heartbeat
// timestamp and the provided threshold duration.
type HealthStatusRequest struct {
	// Nodes is the list of node IDs whose health status
	// needs to be checked.
	Nodes []string `json:"nodes"`

	// LastContactThreshold specifies the maximum allowed
	// duration (in seconds) since the last heartbeat
	// for a node to be considered online.
	LastContactThreshold int `json:"lastContactThreshold"`
}

// HealthStatusDetails represents the health state
// of an individual node.
type HealthStatusDetails struct {
	// Online indicates whether the node is currently online.
	Online bool `json:"online"`

	// LastHeartbeat represents the time (in seconds)
	// since the node last sent a heartbeat.
	LastHeartbeat int `json:"lastHeartbeat"`
}

// HealthStatusResponse represents the response returned
// by the Health Status API.
type HealthStatusResponse struct {
	// BaseResponse contains common response fields such as
	// Success, Error, and ReasonCode.
	common.BaseResponse

	// Data maps node IDs to their corresponding
	// health status details.
	Data map[string]HealthStatusDetails `json:"data"`
}

// GetHealthStatus checks whether one or more nodes are online
// based on their last heartbeat timestamp.
//
// Steps performed by this method:
//  1. Validate the request payload and threshold constraints.
//  2. Marshal the request into JSON format.
//  3. Build and send a POST request to the Health Status API.
//  4. Decode the API response into HealthStatusResponse.
//  5. Map API-level errors into structured SDK errors.
//
// Parameters:
//   - ctx: Context used to control request lifecycle,
//     cancellation, and deadlines.
//   - req: Pointer to HealthStatusRequest containing node IDs
//     and heartbeat threshold.
//
// Returns:
//   - (*HealthStatusResponse, nil) if health status is fetched successfully.
//   - (nil, error) for validation, network, or API-level failures.
func (hm *HealthManagement) GetHealthStatus(
	ctx context.Context,
	req *HealthStatusRequest,
) (*HealthStatusResponse, error) {

	// check if request is nil
	if req == nil {
		return nil, &errors.AnedyaError{
			Message: "health status request cannot be nil",
			Err:     errors.ErrRequestNil,
		}
	}

	// at least one node must be provided
	if len(req.Nodes) == 0 {
		return nil, &errors.AnedyaError{
			Message: "at least one node must be provided",
			Err:     errors.ErrNodesEmpty,
		}
	}

	// validate heartbeat threshold
	if req.LastContactThreshold <= 0 {
		return nil, &errors.AnedyaError{
			Message: "lastContactThreshold must be greater than zero",
			Err:     errors.ErrInvalidTimeRange,
		}
	}

	// threshold must not exceed allowed maximum
	if req.LastContactThreshold > maxHealthThresholdSec {
		return nil, &errors.AnedyaError{
			Message: "lastContactThreshold cannot exceed 7 days",
			Err:     errors.ErrHealthLimitExceeded,
		}
	}

	// build API URL
	url := fmt.Sprintf("%s/v1/health/status", hm.baseURL)

	// convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode health status request",
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
			Message: "failed to build health status request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// send HTTP request
	resp, err := hm.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute health status request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// decode API response
	var apiResp HealthStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode health status response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// handle HTTP-level errors
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// success
	return &HealthStatusResponse{
		Data: apiResp.Data,
	}, nil
}
