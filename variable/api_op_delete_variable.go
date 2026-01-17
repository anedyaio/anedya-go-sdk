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

// DeleteVariableRequest represents the payload sent to the
// Delete Variable API endpoint.
//
// It identifies the variable to be deleted using its unique
// variable key or path.
type DeleteVariableRequest struct {

	// Variable specifies the unique variable key or path
	// that identifies the variable to be deleted.
	//
	// This field is required.
	Variable string `json:"variable"`
}

// DeleteVariableResponse represents the response returned by
// the Delete Variable API endpoint.
//
// It embeds BaseResponse, which contains the standard API
// success flag, error message, and reason code.
type DeleteVariableResponse struct {
	BaseResponse
}

// DeleteVariable deletes an existing variable from the Anedya platform.
//
// The variable to be deleted is identified by its variable key or path,
// which must be provided as a non-empty string.
//
// The method performs the following steps:
//
//  1. Validates the input variable identifier.
//  2. Encodes the request payload as JSON.
//  3. Builds and sends an HTTP request.
//  4. Reads and decodes the API response.
//  5. Maps API errors into structured SDK errors.
//
// Validation errors are returned as sentinel errors defined in the
// errors package. All other failures return *errors.AnedyaError.
func (v *VariableManagement) DeleteVariable(ctx context.Context, variable string) error {

	// 1. Validate input
	if variable == "" {
		return errors.ErrVariableRequired
	}

	// 2. Prepare request payload
	reqPayload := DeleteVariableRequest{
		Variable: variable,
	}

	requestBody, err := json.Marshal(reqPayload)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode DeleteVariable request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/variables/delete", v.baseURL)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build DeleteVariable request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute DeleteVariable request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read DeleteVariable response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode API response
	var apiResp DeleteVariableResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode DeleteVariable response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle HTTP-level errors
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Handle API-level errors
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
