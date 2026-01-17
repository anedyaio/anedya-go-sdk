// Package variable provides APIs to manage variables in the Anedya platform.
package variable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// ListAllVariableRequest represents the payload sent to the
// List Variables API endpoint.
//
// It specifies pagination parameters used to fetch variables
// in batches.
type ListAllVariableRequest struct {

	// Limit specifies the maximum number of variables
	// to return in a single request.
	Limit int `json:"limit"`

	// OffSet specifies the number of variables to skip
	// before starting to return results.
	OffSet int `json:"offset"`
}

// VariableListItem represents a single variable entry returned
// by the List Variables API.
type VariableListItem struct {

	// VariableID is the unique identifier of the variable.
	VariableID string `json:"variableId"`

	// Name is the human-readable name of the variable.
	Name string `json:"name"`

	// Variable is the internal variable key or path.
	Variable string `json:"variable"`

	// TTL specifies the optional time-to-live for the variable.
	TTL int `json:"ttl"`

	// Type specifies the data type of the variable.
	Type string `json:"type"`

	// Description provides additional information about the variable.
	Description string `json:"desc"`
}

// ListAllVariableResponse represents the response returned by
// the List Variables API endpoint.
type ListAllVariableResponse struct {
	BaseResponse

	// CurrentCount indicates the number of variables returned
	// in the current response.
	CurrentCount int `json:"currentCount"`

	// OffSet indicates the starting index used for this response.
	OffSet int `json:"offset"`

	// TotalCount indicates the total number of variables available.
	TotalCount int `json:"totalCount"`

	// NodeParams contains the list of variables returned by the API.
	NodeParams []VariableListItem `json:"nodeParams"`
}

// ListVariablesResult represents the SDK-friendly result returned
// to the caller after converting API response objects.
type ListVariablesResult struct {

	// Variables contains the list of variables returned.
	Variables []Variable

	// CurrentCount indicates the number of variables in this page.
	CurrentCount int

	// OffSet indicates the pagination offset used.
	OffSet int

	// TotalCount indicates the total number of variables available.
	TotalCount int
}

// ListAllVariable retrieves a paginated list of variables from
// the Anedya platform.
//
// Input:
//
//   - limit specifies the maximum number of variables to return in a single request.
//     If set to zero or a negative value, a default value of 100 is used.
//
//   - offset specifies the number of variables to skip before starting to return results.
//     If set to a negative value, it is treated as zero.
//
// Output:
//
//   - On success, the method returns a *ListVariablesResult containing:
//     1) a slice of Variable objects
//     2) pagination metadata (current count, offset, total count)
//
//   - On failure, the method returns a non-nil error.
//
// The method performs the following steps:
//
//   - Validates and normalizes pagination inputs.
//   - Encodes the request payload as JSON.
//   - Builds and sends an HTTP request.
//   - Reads and decodes the API response.
//   - Converts API objects into SDK Variable structures.
//
// Validation and transport failures are returned as structured SDK errors
// using *errors.AnedyaError. API-level failures are mapped using the
// error reason codes returned by the server.
func (v *VariableManagement) ListAllVariable(ctx context.Context, limit int, offset int) (*ListVariablesResult, error) {

	// 1. Validate and normalize inputs
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// 2. Prepare request payload
	reqPayload := ListAllVariableRequest{
		Limit:  limit,
		OffSet: offset,
	}

	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode ListAllVariable request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/variables/list", v.baseURL)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build ListAllVariable request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute ListAllVariable request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read ListAllVariable response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode API response
	var apiResp ListAllVariableResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode ListAllVariable response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle HTTP-level errors
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 9. Convert API response objects to SDK variables
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
