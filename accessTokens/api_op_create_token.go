package accesstokens

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

type Token struct {
	tokenManagement *AccessTokenManagement

	Policy  Policy `json:"policy"`
	TTLSec  int    `json:"ttlSec"`
	TokenID string `json:"tokenId"`
	Token   string `json:"token"`
}

// AccessTokenManagement is the client for access-token APIs.
type AccessTokenManagement struct {
	httpClient *http.Client
	baseURL    string
}

// NewAccessTokenManagement creates a new access-token management client.
func NewAccessTokenManagement(c *http.Client, baseURL string) *AccessTokenManagement {
	return &AccessTokenManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// Policy represents the acess policy for the token
type Policy struct {
	Resources map[string]interface{} `json:"resources,omitempty"`
	Allow     []Permission           `json:"allow,omitempty"`
}

// CreateNewAccessTokenRequest represents the request from CreateNewAccessToken API

type CreateNewAccessTokenRequest struct {
	TTLSec int    `json:"ttlSec"` // Required  1 <= expiry <=7776000
	Policy Policy `json:"policy"` // Required
}

// BaseResponse contains common fields returned by the API.
type BaseResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode"`
}

// CreateNewAccessTokenResponse represents the response from CreateNewAccessToken API
type CreateNewAccessTokenResponse struct {
	BaseResponse
	TokenID string `json:"tokenId"`
	Token   string `json:"token"`
}

// Permission represents an access permission allowed in a token policy
type Permission string

const (
	// Data permissions
	PermissionDataGetSnapshot   Permission = "data::getsnapshot"
	PermissionDataGetLatest     Permission = "data::getlatest"
	PermissionDataGetHistorical Permission = "data::gethistorical"

	// Command persmissions
	PermissionCmdSendCommand  Permission = "cmd::sendcommand"
	PermissionCmdListCommands Permission = "cmd::listcommands"
	PermissionCmdGetStatus    Permission = "cmd::getstatus"
	PermissionCmdInvalidate   Permission = "cmd::invalidate"

	// Variable store permissions
	PermissionVSGetValue   Permission = "vs::getvalue"
	PermissionVSSetValue   Permission = "vs::setvalue"
	PermissionVSScanKeys   Permission = "vs::scankeys"
	PermissionVSDeleteKeys Permission = "vs::deletekeys"

	// Stream permissions
	PermissionStreamsConnect Permission = "streams::connect"

	// Health permissions
	PermissionHealthGetHBStats Permission = "health::gethbstats"
	PermissionHealthGetStatus  Permission = "health::getstatus"
)

func isValidPermission(p Permission) bool {
	switch p {
	case
		PermissionDataGetSnapshot,
		PermissionDataGetLatest,
		PermissionDataGetHistorical,
		PermissionCmdSendCommand,
		PermissionCmdListCommands,
		PermissionCmdGetStatus,
		PermissionCmdInvalidate,
		PermissionVSGetValue,
		PermissionVSSetValue,
		PermissionVSScanKeys,
		PermissionVSDeleteKeys,
		PermissionStreamsConnect,
		PermissionHealthGetHBStats,
		PermissionHealthGetStatus:
		return true
	default:
		return false
	}

}

// CreateNewAccessToken creates a new access token.
func (t *AccessTokenManagement) CreateNewAccessToken(ctx context.Context, input *CreateNewAccessTokenRequest) (*Token, error) {
	// 1. Validate inputs
	if input == nil {
		return nil, &errors.AnedyaError{
			Message: "Input is required",
			Err:     errors.ErrInputRequired,
		}
	}

	const (
		expiryMaxValue int = 7776000 // Maximum expiry in seconds (3 months)
		expiryMinValue int = 1       // Minimum expiry in seconds
	)
	// Validate expiry range
	if input.TTLSec < expiryMinValue || input.TTLSec > expiryMaxValue {
		return nil, &errors.AnedyaError{
			Message: fmt.Sprintf("ttlSec must be between %d and %d seconds", expiryMinValue, expiryMaxValue),
			Err:     errors.ErrExpiryRequried,
		}
	}
	// Validate policy
	if len(input.Policy.Allow) == 0 {
		return nil, &errors.AnedyaError{
			Message: "must contain at least one permission",
			Err:     errors.ErrPolicyRequired,
		}
	}

	for _, p := range input.Policy.Allow {
		if !isValidPermission(p) {
			return nil, &errors.AnedyaError{
				Message: "Invalid permission",
				Err:     errors.ErrInavalidPermission,
			}
		}
	}
	// 2. Prepare payload
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode create token request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Create HTTP Request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
	url := fmt.Sprintf("%s/v1/access/tokens/create", t.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build create token request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute Request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute create token request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response data
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read create token response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 7. Decode response
	var apiResp CreateNewAccessTokenResponse
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode create token response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 6. Check for status codes
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return &Token{
		tokenManagement: t,
		TokenID:         apiResp.TokenID,
		Token:           apiResp.Token,
		Policy:          input.Policy,
		TTLSec:          input.TTLSec,
	}, nil
}
