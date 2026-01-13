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

// ListAllVariableRequest represents the request body for ListAllVariable API
type ListAllVariableRequest struct {
	Limit  int `json:"limit"`
	OffSet int `json:"offset"`
}
type VariableListItem struct {
	VariableID  string `json:"variableId"`
	Name        string `json:"name"`
	Variable    string `json:"variable"`
	TTL         int    `json:"ttl"`
	Type        string `json:"type"`
	Description string `json:"desc"`
}

// ListAllVariableResponse represents the response from ListAllVariable API
type ListAllVariableResponse struct {
	BaseResponse
	CurrentCount int                `json:"currentCount"`
	OffSet       int                `json:"offset"`
	TotalCount   int                `json:"totalCount"`
	NodeParams   []VariableListItem `json:"nodeParams"`
}

// ListVariablesResult represents the paginated result with variables
type ListVariablesResult struct {
	Variables    []Variable
	CurrentCount int
	OffSet       int
	TotalCount   int
}

// ============================================
// API Methods
// ============================================

// ListAllVariable lists all variables
func (v *VariableManagement) ListAllVariable(ctx context.Context, limit, offset int) (*ListVariablesResult, error) {
	// 1. Validating inputs
	if limit == 0 {
		limit = 100
	}
	if offset == 0 {
		offset = 0
	}

	// 2. Preparing payload
	reqPayload := ListAllVariableRequest{
		Limit:  limit,
		OffSet: offset,
	}
	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	// 3. Create Request
	// assuming baseUrl to be "https://api.ap-in-1.anedya.io"
	url := fmt.Sprintf("%s/v1/variables/list", v.baseURL)
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
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 6. Check for status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("api failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 7. Decode response
	var apiResp ListAllVariableResponse

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	// 8. Convert VariableListItem into Variable
	variables := make([]Variable, len(apiResp.NodeParams))
	for i, item := range apiResp.NodeParams {
		variables[i] = Variable{
			variableManagement: v,
			VariableID:         item.VariableID,
			Name:               item.Name,
			Type:               item.Type,
			Variable:           item.Variable,
			Description:        item.Description,
			TTL:                item.TTL,
		}
	}

	return &ListVariablesResult{
		Variables:    variables,
		CurrentCount: apiResp.CurrentCount,
		OffSet:       apiResp.OffSet,
		TotalCount:   apiResp.TotalCount,
	}, nil
}
