package variable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Variable struct {
	variableManagement *VariableManagement

	VariableID  string `json:"variableId,omitempty"`
	Type        string `json:"type"` // Required (Possible values: [float, geo])
	Name        string `json:"name"` // Required
	Description string `json:"desc,omitempty"`
	Variable    string `json:"variable"` // Required
	TTL         int    `json:"ttl,omitempty"`
}

// Variable Management Object
type VariableManagement struct {
	httpClient *http.Client
	baseURL    string
}

// ============================================
// Constructor
// ============================================
func NewVariableManagement(c *http.Client, baseURL string) *VariableManagement {
	return &VariableManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// ============================================
// API request and response Structs
// ============================================

// CreateVariableRequest represents the request body for CreateVariable API
type CreateVariableRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"desc,omitempty"`
	Variable    string `json:"variable"`
	TTL         int    `json:"ttl,omitempty"`
}
type BaseResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode"`
}

type CreateVariableResponse struct {
	BaseResponse
	VariableID string `json:"variableId"`
}

// ValidVariableTypes contains all valid variable types
var ValidVariableTypes = []string{"geo", "float"}      // remove global variable

func isValidVariableType(variableType string) bool {
	for _, v := range ValidVariableTypes {
		if strings.EqualFold(variableType, v) {
			return true
		}
	}
	return false
}

// ============================================
// API Methods
// ============================================

// Create a new variable method
func (v *VariableManagement) CreateVariable(ctx context.Context, variable *Variable) (*Variable, error) {
	// 1. Validating Inputs
	if variable.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if variable.Variable == "" {
		return nil, fmt.Errorf("variable is required")
	}
	// Validate type using the array
	if !isValidVariableType(variable.Type) {
		return nil, fmt.Errorf("type must be one of: %s", strings.Join(ValidVariableTypes, ", "))
	}

	// 2. Preparing payload
	reqPayload := CreateVariableRequest{
		Type:        variable.Type,
		Name:        variable.Name,
		Description: variable.Description,
		Variable:    variable.Variable,
		TTL:         variable.TTL,
	}
	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	// 3. Create Request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
	url := fmt.Sprintf("%s/v1/variables/create", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute Request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 5. Read response data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 6. Check for status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("api failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 7. Decode response
	var apiResp CreateVariableResponse

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return &Variable{
		VariableID: apiResp.VariableID,
	}, nil
}
