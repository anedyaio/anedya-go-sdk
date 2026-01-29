// Package valuestore provides APIs to store and manage key-value data
// in the Anedya platform.
//
// The Value Store allows applications to persist data at either a
// global project scope or a specific node scope, with support for
// multiple data types such as string, boolean, float, and binary.
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

// Value represents a stored key-value entry in the Anedya Value Store.
//
// This structure is typically returned by read APIs or used internally
// to represent persisted data.
type Value struct {

	// valueStoreManagement holds the internal client used to
	// perform operations on this value.
	//
	// This field is not serialized and is used internally by the SDK.
	valueStoreManagement *ValueStoreManagement

	// NameSpace specifies the scope and identifier for the value.
	NameSpace NameSpace `json:"namespace"`

	// Key uniquely identifies the value within the namespace.
	Key string `json:"key"`

	// Value contains the actual stored data.
	//
	// Since the Value Store supports multiple data types, this field holds
	// the value as a generic interface{}. The underlying Go type depends
	Value json.RawMessage `json:"value"`

	// Type indicates the data type of the stored value.
	Type ValueType `json:"type"`

	// Size represents the size of the stored value in bytes.
	Size int `json:"size"`
}

// SetValueRequest represents the payload sent to the
// Create Variable API endpoint.
//
// All required fields must be provided before calling SetValueRequest.
type SetValueRequest struct {

	// NameSpace specifies about the key-value pair that is stored.
	//
	// It defines the storage scope (global or node) and the
	// associated identifier used to locate the value.
	NameSpace NameSpace `json:"namespace"`

	// Key is the unique identifier for the value.
	Key string `json:"key"`

	// Value contains the data to be stored.
	Value string `json:"value"`

	// Type specifies the data type of the value.
	Type ValueType `json:"type"`
}

// SetValueRespone represents the API response for SetValue.
type SetValueResponse struct {

	// BaseResponse contains common API response fields such as
	// success flag, error message, and machine-readable reason code.
	common.BaseResponse
}

// SetValue stores or updates a value in the Anedya Value Store.
//
// The value can be stored either globally or under a specific node,
// based on the provided namespace configuration.
//
// Input Parameters:
//   - ctx: Context used for request cancellation and deadlines.
//   - input: Pointer to SetValueRequest containing value details.
//
// The method performs the following steps:
//
//  1. Validates request fields and enums.
//  2. Encodes the request payload as JSON.
//  3. Builds an HTTP POST request.
//  4. Executes the request.
//  5. Reads and decodes the API response.
//  6. Maps API error codes into SDK structured errors.
//
// Validation errors return sentinel errors wrapped inside *errors.AnedyaError.
// API failures are converted using errors.GetError().
//
// On success, the method returns nil.
func (v *ValueStoreManagement) SetValue(ctx context.Context, input *SetValueRequest) error {

	// 1. Validate input
	if input == nil {
		return &errors.AnedyaError{
			Message: "input is required",
			Err:     errors.ErrInputRequired,
		}
	}

	// Validate scope
	// Check if Namespace Scope is empty (Namespace is Required)
	if input.NameSpace.Scope == "" {
		return &errors.AnedyaError{
			Message: "namespace scope is required",
			Err:     errors.ErrNamespaceScopeRequired,
		}
	}

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

	// Validate value
	if input.Value == "" {
		return &errors.AnedyaError{
			Message: "value is required",
			Err:     errors.ErrValueRequired,
		}
	}

	// Validate value type
	if !isValidValueType(input.Type) {
		return &errors.AnedyaError{
			Message: "invalid value type",
			Err:     errors.ErrInvalidValueType,
		}
	}

	// 2. Encode request body
	requestBody, err := json.Marshal(input)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to encode set value request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP request
	url := fmt.Sprintf("%s/v1/valuestore/setValue", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to build set value request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to execute set value request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.AnedyaError{
			Message: "failed to read set value response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode response
	var apiResp SetValueResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return &errors.AnedyaError{
			Message: "failed to decode set value response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 8. Handle API-level errors.
	if !apiResp.Success {
		return errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return nil
}
