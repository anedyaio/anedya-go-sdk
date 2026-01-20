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

// RevokeAcessTokenRequest represents the request from RevokeAcsessToken API
type RevokeAccessTokenRequest struct {
	TokenID string `json:"tokenId"`
}

// RevokeAccessTokenResponse represents the response from RevokeAccessToken API
type RevokeAccessTokenResponse struct {
	BaseResponse
}

// RevokeAccessToken revokes any issued token by providing tokenID
func (t *AccessTokenManagement) RevokeAccessToken(ctx context.Context, tokenId string) error {
	// 1. Validate Inputs
	if tokenId == "" {
		return &errors.AnedyaError{
			Message: "tokenId is required",
			Err:     errors.ErrTokenIdRequired,
		}
	}

	// 2. Prepare payload
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

	// 3. Create HTTP request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
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

	// 4. Execute Request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute revoke token request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response data
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read revoke token response ",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 7. Decode response
	var apiResp RevokeAccessTokenResponse

	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode revoke token response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 6. Check for status codes
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
