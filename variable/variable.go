package variable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Variable struct {
	variableManagement *VariableManagement

	VariableID   string `json:"variableId,omitempty"`
	Type         string `json:"type"` // Required (Possible values: [float, geo])
	Name         string `json:"name"` // Required
	Description  string `json:"desc,omitempty"`
	VariableCode string `json:"variable"` // Required
	TTL          int    `json:"ttl,omitempty"`
}

// Variable Management Object
type VariableManagement struct {
	httpClient *http.Client
	baseURL    string
}

// Constructor
func NewVariableManagement(c *http.Client, baseURL string) *VariableManagement {
	return &VariableManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// API Response Structs
type BaseResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode"`
}

type CreateVariableResponse struct {
	BaseResponse
	VariableID string `json:"variableId"`
}

// CreateVariable Method
func (v *VariableManagement) CreateVariable(ctx context.Context, variable *Variable) (*Variable, error) {
	// 1. Validating Inputs
	if variable.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if variable.VariableCode == "" {
		return nil, fmt.Errorf("variable is required")
	}
	if variable.Type != "float" && variable.Type != "geo" {
		return nil, fmt.Errorf("type must be float or geo")
	}

	// 2. Preparing payload
	requestBody, err := json.Marshal(variable)
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api failed with statu %d: %s", resp.StatusCode, string(body))
	}

	// 7. Decode response
	type apiResponse struct {
		Sucess     bool   `json:"success"`
		Error      string `json:"error"`
		ReasonCode string `json:"reasonCode"`
		VariableID string `json:"variableId"`
	}

	var apiResp CreateVariableResponse

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return &Variable{
		variableManagement: v,
		VariableID:         apiResp.VariableID,
		Name:               variable.Name,
		Type:               variable.Type,
		VariableCode:       variable.VariableCode,
		Description:        variable.Description,
		TTL:                variable.TTL,
	}, nil
}
