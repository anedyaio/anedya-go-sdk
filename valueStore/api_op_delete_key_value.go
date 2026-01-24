// Package valuestore provides APIs to store and manage key-value data
// in the Anedya platform.
package valuestore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anedyaio/anedya-go-sdk/common"
	"github.com/anedyaio/anedya-go-sdk/errors"
)

// DeleteKeyValueRequest represents the payload sent to the
// Delete Key-Value API endpoint.
//
// It identifies the key-value pair to be deleted using a namespace
// (global or node scope) and the corresponding key.
type DeleteKeyValueRequest struct {

	// NameSpace specifies about the key-value pair that is stored.
	//
	// It defines the storage scope (global or node) and the
	// associated identifier used to locate the value.
	NameSpace NameSpace `json:"namespace"`

	// Key is the key of the value to be deleted.
	Key string `json:"key,omitempty"`
}

// DeleteKeyValueResponse represents the response returned by the
// Delete Key-Value API endpoint.
type DeleteKeyValueResponse struct {

	// BaseResponse contains common API response fields such as
	// success flag, error message, and machine-readable reason code.
	common.BaseResponse
}

// DeleteKeyValuePair deletes a key-value pair from the Anedya Value Store.
//
// The key-value pair is identified using a namespace (global or node scope)
// and a key provided through the DeleteKeyValueRequest structure.
//
// The method performs the following steps:
//
//  1. Validates the input payload and namespace fields.
//  2. Encodes the request payload as JSON.
//  3. Builds and sends an HTTP request to the API.
//  4. Reads and decodes the API response.
//  5. Converts API errors into structured SDK errors.
//
// Validation failures return sentinel errors defined in the errors package.
// Network and API failures return *errors.AnedyaError.
//
// On success, the method returns nil.
func (v *ValueStoreManagement) DeleteKeyValuePair(ctx context.Context, input *DeleteKeyValueRequest) error {

	// 1. Validate input
	if input == nil {
		return &errors.AnedyaError{
			Message: "input is required",
			Err:     errors.ErrInputRequired,
		}
	}

	// Validate scope
	if !isValidScope(input.NameSpace.Scope) {
		return &errors.AnedyaError{
			Message: "invalid namespace scope",
			Err:     errors.ErrInvalidNamespaceScope,
		}
	}

	// Validate namespace ID
	if input.NameSpace.Id == "" {
		return &errors.AnedyaError{
			Message: "namespace id is required",
			Err:     errors.ErrValueNamespaceIdRequired,
		}
	}

	// Validate key
	if input.Key == "" {
		return &errors.AnedyaError{
			Message: "key is required",
			Err:     errors.ErrValueKeyRequired,
		}
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(input)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode delete key value pair request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/valuestore/delete", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build delete key value pair request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute delete key value pair request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read delete key value pair response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp DeleteKeyValueResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode delete key value pair response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle HTTP-level errors
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Handle API-level errors.
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
