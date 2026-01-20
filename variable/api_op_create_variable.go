// Package variable provides APIs to manage variables in the Anedya platform.
package variable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/anedyaio/anedya-go-sdk/errors"
)

// Variable represents a variable resource returned by the Anedya API.
//
// A Variable uniquely identifies a telemetry or configuration field
// that can be associated with devices or nodes.
//
// The fields correspond to the values returned by the API when a
// variable is created or queried.
type Variable struct {

	// variableManagement holds the internal client used to
	// perform operations on this variable.
	//
	// This field is not serialized and is used internally
	// by the SDK.
	variableManagement *VariableManagement

	// VariableID is the unique identifier assigned by the API
	// to the variable.
	VariableID string `json:"variableId,omitempty"`

	// Type specifies the data type of the variable.
	//
	// Supported values are:
	//   - "float"
	//   - "geo"
	Type string `json:"type"`

	// Name is the human-readable name of the variable.
	Name string `json:"name"`

	// Description provides an optional description
	// of the variable.
	Description string `json:"desc,omitempty"`

	// Variable is the variable key or path used internally
	// by the Anedya platform.
	Variable string `json:"variable"`

	// TTL specifies the optional time-to-live (in seconds)
	// for values associated with this variable.
	TTL int `json:"ttl,omitempty"`
}

// VariableManagement provides methods to create, update, delete,
// and retrieve variables from the Anedya API.
//
// It wraps an HTTP client and a base URL used to construct requests.
type VariableManagement struct {

	// httpClient is the underlying HTTP client used to
	// perform API requests.
	httpClient *http.Client

	// baseURL is the root API endpoint used for all
	// variable management requests.
	baseURL string
}

// NewVariableManagement creates a new VariableManagement client.
//
// The provided http.Client is used for all network communication,
// and baseURL specifies the API server address.
func NewVariableManagement(c *http.Client, baseURL string) *VariableManagement {
	return &VariableManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// CreateVariableRequest represents the payload sent to the
// Create Variable API endpoint.
//
// All required fields must be provided before calling CreateVariable.
type CreateVariableRequest struct {

	// Type specifies the variable data type.
	//
	// Supported values are:
	//   - "float"
	//   - "geo"
	Type string `json:"type"`

	// Name specifies the human-readable name of the variable.
	Name string `json:"name"`

	// Description provides an optional explanation of
	// what the variable represents.
	Description string `json:"desc,omitempty"`

	// Variable specifies the unique variable key or path.
	Variable string `json:"variable"`

	// TTL specifies the optional expiration time (in seconds)
	// for values associated with this variable.
	TTL int `json:"ttl,omitempty"`
}

// BaseResponse represents common fields returned by all
// Anedya API responses.
type BaseResponse struct {

	// Success indicates whether the API request was successful.
	Success bool `json:"success"`

	// Error contains the error message returned by the API
	// when Success is false.
	Error string `json:"error"`

	// ReasonCode contains the machine-readable error code
	// used for SDK error mapping.
	ReasonCode string `json:"reasonCode"`
}

// CreateVariableResponse represents the response returned by
// the Create Variable API endpoint.
type CreateVariableResponse struct {
	BaseResponse

	// VariableID is the identifier of the newly created variable.
	VariableID string `json:"variableId"`
}

// CreateVariable creates a new variable in the Anedya platform.
//
// The request is provided using a *CreateVariableRequest structure, which
// defines the variable name, type, identifier, and optional metadata such
// as description and TTL.
//
// On success, the method returns a *Variable containing the unique
// VariableID assigned by the Anedya API.
//
// The method performs the following steps:
//
//  1. Validates the request payload.
//  2. Encodes the payload as JSON.
//  3. Builds and sends an HTTP request.
//  4. Reads and decodes the API response.
//  5. Maps API errors into structured SDK errors.
//
// Validation errors are returned as sentinel errors defined in the
// errors package. All other failures return *errors.AnedyaError.
func (v *VariableManagement) CreateVariable(ctx context.Context, input *CreateVariableRequest) (*Variable, error) {

	// 1. Validate input payload.
	if input == nil {
		return nil, &errors.AnedyaError{
			Message: "Input is required",
			Err:     errors.ErrInputRequired,
		}
	}

	if input.Name == "" {
		return nil, &errors.AnedyaError{
			Message: "Input name is requried",
			Err:     errors.ErrVariableNameRequired,
		}
	}

	if input.Variable == "" {
		return nil, &errors.AnedyaError{
			Message: "Input variable is required",
			Err:     errors.ErrVariableRequired,
		}
	}

	validTypes := map[string]bool{
		"geo":   true,
		"float": true,
	}
	if !validTypes[strings.ToLower(input.Type)] {
		return nil, errors.ErrVariableTypeRequired
	}

	// 2. Encode request body.
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode CreateVariable request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request.
	url := fmt.Sprintf("%s/v1/variables/create", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build CreateVariable request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute request.
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute CreateVariable request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read CreateVariable response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp CreateVariableResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode CreateVariable response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle HTTP-level errors.
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Handle API-level errors.
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 9. Return created variable.
	return &Variable{
		variableManagement: v,
		VariableID:         apiResp.VariableID,
		Type:               input.Type,
		Name:               input.Name,
		Description:        input.Description,
		Variable:           input.Variable,
		TTL:                input.TTL,
	}, nil
}
