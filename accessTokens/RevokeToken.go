package accesstokens

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ============================================
// API request and response structs
// ============================================

// RevokeAcessTokenRequest represents the request from RevokeAcsessToken API
/*
	Request Body example:
	{
	"tokenId": "string"
	}
*/
type RevokeAccessTokenRequest struct {
	TokenID string `json:"tokenId"`
}

// RevokeAccessTokenResponse represents the response from RevokeAccessToken API
type RevokeAccessTokenResponse struct {
	BaseResponse
}

// ============================================
// API Methods
// ============================================
// RevokeAccessToken revokes any issued token by providing tokenID
func (t *AccessTokenManagement) RevokeAccessToken(ctx context.Context, tokenId string) error {
	// 1. Validate Inputs
	if tokenId == "" {
		return fmt.Errorf("tokenId is required")
	}

	// 2. Prepare payload
	reqPayload := RevokeAccessTokenRequest{
		TokenID: tokenId,
	}
	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request %w", err)
	}

	// 3. Create HTTP request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
	url := fmt.Sprintf("%s/v1/access/tokens/revoke", t.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute Request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 5. Read response data
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 6. Check for status codes
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("api failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// 7. Decode response
	var apiResp RevokeAccessTokenResponse

	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Error)
	}

	return nil
}
