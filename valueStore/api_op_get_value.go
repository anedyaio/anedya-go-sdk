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

// GetValueRequest represents the payload sent to the
// Get Value API endpoint.
//
// It specifies the namespace and key for which the stored
// value should be retrieved.
type GetValueRequest struct {

	// NameSpace defines the logical scope in which the key exists.
	//
	// It contains:
	//   - Scope: Determines whether the value is stored at project level ("global")
	//            or under a specific node ("node").
	//   - Id:    A project-wide unique identifier for "global" scope, or a valid
	//            node ID when the scope is "node".
	NameSpace NameSpace `json:"namespace,omitempty"`

	// Key is the name of the key whose value should be fetched
	// from the value store.
	Key string `json:"key,omitempty"`
}

// GetValueResponse represents the response returned by the
// Get Value API endpoint.
//
// It contains both the common API response fields and the
// actual stored key-value data.
type GetValueResponse struct {
	common.BaseResponse

	// NameSpace indicates the namespace from which the value
	// was retrieved, including its scope and identifier.
	NameSpace NameSpace `json:"namespace"`

	// Key is the name of the stored key returned by the API.
	Key string `json:"key"`

	// Type specifies the data type of the stored value.
	//
	// Possible values:
	//   - "string"
	//   - "binary"
	//   - "float"
	//   - "boolean"
	Type ValueType `json:"type"`

	// Size represents the size of the stored value in bytes.
	Size int `json:"size"`

	// Value contains the actual stored value.
	//
	// The concrete Go type depends on the Type field:
	//   - string   → string
	//   - boolean  → bool
	//   - float    → float64
	//   - binary   → string (base64 encoded)
	Value interface{} `json:"value"`
}

// GetValue retrieves a stored value from the Anedya value store.
//
// The method fetches a key-value pair from the specified namespace
// (global or node-level) and returns its metadata and actual value.
//
// Parameters:
//   - ctx:   Context used to control request lifetime and cancellation.
//   - input: Pointer to GetValueRequest containing the namespace and key.
//
// Returns:
//   - *GetValueResponse: On success, contains the stored value along with its
//     namespace, key name, data type, and size.
//   - error: Returns a structured *errors.AnedyaError if the request fails due to
//     validation issues, network errors, API errors, or response decoding failures.
//
// The method performs the following steps:
//
//  1. Validates the request payload.
//  2. Encodes the request body as JSON.
//  3. Builds the HTTP request.
//  4. Sends the request to the Anedya API.
//  5. Reads and decodes the response body.
//  6. Handles HTTP-level and API-level errors.
//  7. Returns the parsed response on success.
func (v *ValueStoreManagement) GetValue(ctx context.Context, input *GetValueRequest) (*Value, error) {

	// 1. Validate Input Struct
	if input == nil {
		return nil, &errors.AnedyaError{
			Message: "get value request cannot be nil",
			Err:     errors.ErrInputRequired,
		}
	}

	// Validate Namespace Id
	if input.NameSpace.Id == "" {
		return nil, &errors.AnedyaError{
			Message: "namespace id is required",
			Err:     errors.ErrValueNamespaceIdRequired,
		}
	}

	// Validate scope
	if !isValidScope(input.NameSpace.Scope) {
		return nil, &errors.AnedyaError{
			Message: "invalid namespace scope",
			Err:     errors.ErrInvalidNamespaceScope,
		}
	}

	// Validate key
	if input.Key == "" {
		return nil, &errors.AnedyaError{
			Message: "key is required",
			Err:     errors.ErrValueKeyRequired,
		}
	}

	// 2. Encode Request
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode get value request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/valuestore/getValue", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build get value request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 4. Execute request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute get value request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read get value response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp GetValueResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode get value response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 7. Handle HTTP-level errors
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	// 8. Handle API-level errors.
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return &Value{
		NameSpace: input.NameSpace,
		Key:       apiResp.Key,
		Value:     apiResp.Value,
		Type:      apiResp.Type,
		Size:      apiResp.Size,
	}, nil
}
