package accessTokens

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Token struct {
	tokenManagement *AccessTokenManagement

	Policy
	Expiry  int    `json:"expiry"`
	TokenID string `json:"tokenId"`
	Token   string `json:"token"`
}

// AccessTokenManagement is the client for access-token APIs.
type AccessTokenManagement struct {
	httpClient *http.Client
	baseURL    string
}

// ============================================
// Constructor
// ============================================

// NewAccessTokenManagement creates a new access-token management client.
func NewAccessTokenManagement(c *http.Client, baseURL string) *AccessTokenManagement {
	return &AccessTokenManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// ============================================
// API request and response structs
// ============================================

// Policy represents the acess policy for the token
type Policy struct {
	Resources map[string]struct{} `json:"resources"`
	Allow     []string            `json:"allow"`
}

// BaseResponse contains common fields returned by the API.
type BaseResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode"`
}

// CreateNewAccessTokenRequest represents the request from CreateNewAccessToken API
/*
	Request Body example:
	{
	"expiry": 0,
	"policy": {
		"resources": {},
		"allow": [
		"string"
		]
	}
	}
*/
type CreateNewAccessTokenRequest struct {
	Expiry int `json:"expiry"` // Required  1 <= expiry <=7776000
	Policy     // Required
}

// CreateNewAccessTokenResponse represents the response from CreateNewAccessToken API
type CreateNewAccessTokenResponse struct {
	BaseResponse
	TokenID string `json:"tokenId"`
	Token   string `json:"token"`
}

const (
	expiryMaxValue int = 7776000 // Maximum expiry in seconds (3 months)
	expiryMinValue int = 1       // Minimum expiry in seconds
)

// CreateNewAccessToken creates a new access token.
func (t *AccessTokenManagement) CreateNewAccessToken(ctx context.Context, token *CreateNewAccessTokenRequest) (*Token, error) {
	// 1. Validate inputs
	if token == nil {
		return nil, fmt.Errorf("request body is empty")
	}
	// Validate expiry range
	if token.Expiry < expiryMinValue || token.Expiry > expiryMaxValue {
		return nil, fmt.Errorf("Expiry value must lie between %d and %d", expiryMinValue, expiryMaxValue)
	}
	// Validate policy
	if len(token.Policy.Allow) == 0 {
		return nil, fmt.Errorf("policy.allow must contain at least one permission")
	}

	// 2. Prepare payload
	reqPayload := CreateNewAccessTokenRequest{
		Expiry: token.Expiry,
		Policy: token.Policy,
	}
	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	// 3. Create HTTP Request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
	url := fmt.Sprintf("%s/v1/access/tokens/create", t.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute Request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 5. Read response data
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 6. Check for status codes
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("api failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// 7. Decode response
	var apiResp CreateNewAccessTokenResponse

	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return &Token{
		TokenID: apiResp.TokenID,
		Token:   apiResp.Token,
	}, nil
}
