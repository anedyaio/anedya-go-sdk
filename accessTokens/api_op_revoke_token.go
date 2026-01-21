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

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// RevokeAccessTokenRequest represents the payload sent to the
// Revoke Access Token API endpoint.
//
// It identifies the token to be revoked using its unique token ID.
type RevokeAccessTokenRequest struct {

	// TokenID is the unique identifier of the access token
	// that should be revoked.
	//
	// This field is required.
	TokenID string `json:"tokenId"`
}

// RevokeAccessTokenResponse represents the response returned by
// the Revoke Access Token API endpoint.
//
// It embeds BaseResponse, which contains the standard API
// success flag, error message, and reason code.
type RevokeAccessTokenResponse struct {
	BaseResponse
}

// RevokeAccessToken revokes an existing access token in the Anedya platform.
//
// Input:
//   - ctx: request context
//   - input: RevokeAccessTokenRequest containing TokenID
//
// Output:
//   - error on failure
//
// The method performs the following steps:
//
//  1. Validates the input token identifier.
//  2. Encodes the request payload as JSON.
//  3. Builds and sends an HTTP request.
//  4. Reads and decodes the API response.
//  5. Maps API errors into structured SDK errors.
//
// Validation errors are returned as sentinel errors defined in the
// errors package. All other failures return *errors.AnedyaError.
func (t *AccessTokenManagement) RevokeAccessToken(ctx context.Context, tokenId string) error {
	// Step 1: Validate the input token identifier.
	if tokenId == "" {
		return &errors.AnedyaError{
			Message: "tokenId is required",
			Err:     errors.ErrTokenIdRequired,
		}
	}

	// Step 2: Construct the request payload.
	reqPayload := RevokeAccessTokenRequest{
		TokenID: tokenId,
	}
	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode revoke token request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// Step 3: Build the HTTP request.
	url := fmt.Sprintf("%s/v1/access/tokens/revoke", t.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build revoke token request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Step 4: Execute the HTTP request.
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute revoke token request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Step 5: Read the response body.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read revoke token response ",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// Step 6: Decode the API response.
	var apiResp RevokeAccessTokenResponse

	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode revoke token response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// Step 7: Handle HTTP-level errors.
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Step 8: Handle API-level errors.
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// Step 9: Token successfully revoked.
	return nil
}
