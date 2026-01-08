package variable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ============================================
// API request and response Structs
// ============================================

// DeleteVariableRequest represents the request body for DeleteVariable API
type DeleteVariableRequest struct {
	Variable string `json:"variable"` // Required
}

// DeleteVariableResponse represents the response from DeleteVariable API
type DeleteVariableResponse struct {
	BaseResponse
}

// ============================================
// API Methods
// ============================================

// DeleteVariable deletes a variable
func (v *VariableManagement) DeleteVariable(ctx context.Context, variable string) error {
	// 1. Validating inputs
	if variable == "" {
		return fmt.Errorf("variable is required")
	}

	// 2. Preparing payload
	reqPayload := DeleteVariableRequest{
		Variable: variable,
	}
	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}

	// 3. Create Request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
	url := fmt.Sprintf("%s/v1/variables/delete", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute Request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 5. Read response data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 6. Check for status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("api failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 7. Decode response
	var apiResp DeleteVariableResponse

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return fmt.Errorf("failed to parse reponse: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Error)
	}

	return nil
}
