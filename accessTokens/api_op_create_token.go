// Package accesstokens provides APIs to manage access tokens
// in the Anedya platform.
package accesstokens

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

// Token represents an access token returned by the Anedya API.
//
// A Token allows authenticated access to platform resources
// based on the attached access policy.
type Token struct {

	// tokenManagement holds the internal client used to
	// perform operations related to this token.
	//
	// This field is not serialized and is used internally
	// by the SDK.
	tokenManagement *AccessTokenManagement

	// Policy defines the access rules associated with the token,
	// including allowed resources and permissions.
	Policy Policy `json:"policy"`

	// TTLSec specifies the time-to-live of the token in seconds.
	//
	// After this duration, the token expires and can no longer
	// be used for authentication.
	TTLSec int `json:"ttlSec"`

	// TokenID is the unique identifier assigned by the API
	// to the access token.
	TokenID string `json:"tokenId"`

	// Token is the actual secret value used for authentication
	// in API requests.
	Token string `json:"token"`
}

// Policy represents the access policy attached to an access token.
//
// It defines which resources can be accessed and what actions
// are allowed on those resources.
type Policy struct {

	// Resources specifies the set of resources that the token
	// can access.
	//
	// The structure of this object depends on the Anedya Access
	// Policy specification (for example: nodes, devices, etc.).
	Resources map[string]interface{} `json:"resources,omitempty"`

	// Allow lists the permissions granted to the token.
	//
	// Each permission must be a valid Permission constant
	// defined by the SDK.
	Allow []Permission `json:"allow,omitempty"`
}

// CreateNewAccessTokenRequest represents the payload sent to the
// Create Access Token API endpoint.
//
// All required fields must be provided before calling CreateNewAccessToken.
type CreateNewAccessTokenRequest struct {

	// TTLSec specifies the token expiration time in seconds.
	//
	// Valid range:
	//   1 <= TTLSec <= 7776000 (3 months)
	TTLSec int `json:"ttlSec"`

	// Policy specifies the access policy that controls
	// what the token is allowed to do.
	Policy Policy `json:"policy"`
}

// CreateNewAccessTokenResponse represents the response returned by
// the Create Access Token API endpoint.
type CreateNewAccessTokenResponse struct {
	common.BaseResponse

	// TokenID is the identifier of the newly created token.
	TokenID string `json:"tokenId"`

	// Token is the generated secret token value.
	Token string `json:"token"`
}

// Permission represents a single access permission
// that can be granted to an access token.
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

// isValidPermission validates whether the given permission
// is supported by the SDK.
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

// CreateNewAccessToken creates a new access token in the Anedya platform.
//
// Input:
//   - ctx: request context
//   - input: CreateNewAccessTokenRequest containing TTL and policy
//
// Output:
//   - *Token on success
//   - error on failure
//
// The method performs the following steps:
//
//  1. Validates the input payload.
//  2. Encodes the payload as JSON.
//  3. Builds and sends an HTTP request.
//  4. Reads and decodes the API response.
//  5. Maps API errors into structured SDK errors.
//
// Validation errors are returned as sentinel errors defined in the
// errors package. All other failures return *errors.AnedyaError.
func (t *AccessTokenManagement) CreateNewAccessToken(ctx context.Context, input *CreateNewAccessTokenRequest) (*Token, error) {
	// Validate the input request
	if input == nil {
		return nil, &errors.AnedyaError{
			Message: "Input is required",
			Err:     errors.ErrInputRequired,
		}
	}

	// Define valid bounds for token expiry (in seconds).
	const (
		expiryMaxValue int = 7776000 // Maximum expiry in seconds (3 months)
		expiryMinValue int = 1       // Minimum expiry in seconds
	)

	// Validate that the provided TTL is within the allowed range.
	if input.TTLSec < expiryMinValue || input.TTLSec > expiryMaxValue {
		return nil, &errors.AnedyaError{
			Message: fmt.Sprintf("ttlSec must be between %d and %d seconds", expiryMinValue, expiryMaxValue),
			Err:     errors.ErrExpiryRequried,
		}
	}

	// Ensure that at least one permission is defined in the policy.
	if len(input.Policy.Allow) == 0 {
		return nil, &errors.AnedyaError{
			Message: "must contain at least one permission",
			Err:     errors.ErrPolicyRequired,
		}
	}

	// Validate that all provided permissions are supported by the SDK.
	for _, p := range input.Policy.Allow {
		if !isValidPermission(p) {
			return nil, &errors.AnedyaError{
				Message: "Invalid permission",
				Err:     errors.ErrInavalidPermission,
			}
		}
	}

	// Encode the request payload
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode create token request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Construct the HTTP request for the API endpoint.
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

	// Send the HTTP request to the API server
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute create token request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Read the raw response body from the API.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read create token response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// Decode the JSON response into the response structure.
	var apiResp CreateNewAccessTokenResponse
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode create token response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Handle HTTP-level errors.
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Handle API-level errors.
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Construct and return the SDK Token object.
	return &Token{
		tokenManagement: t,
		TokenID:         apiResp.TokenID,
		Token:           apiResp.Token,
		Policy:          input.Policy,
		TTLSec:          input.TTLSec,
	}, nil
}
