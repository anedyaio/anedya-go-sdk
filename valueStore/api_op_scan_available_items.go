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

// Filter specifies the criteria used to narrow down the scan results.
type Filter struct {
	// NameSpace defines the scope to search within.
	// This field is required.
	NameSpace NameSpace `json:"namespace"`
}

// ListValuesRequest represents the input parameters for the List Values API.
type ListValuesRequest struct {

	// Filter contains the scoping criteria for the request (e.g., Namespace)
	Filter Filter `json:"filter"`

	// OrderBy specifies which field to sort the results by.
	// This field is required.
	// Valid values: "namespace", "key", "created".
	OrderBy ScanOrderBy `json:"orderby"`

	// Order specifies the direction of the sort (ascending or descending).
	// This field is optional.
	// Valid values: "asc", "desc".
	Order SortOrder `json:"order"`

	// Limit specifies the maximum number of items to return in one page.
	// If 0, the server applies its default limit.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of items to skip before starting to collect the result set.
	// Used for pagination.
	Offset int `json:"offset,omitempty"`
}

// ValueItem represents the metadata of a single key found in the Value Store.
type ValueItem struct {
	// NameSpace indicates where this key is stored (scope and ID).
	NameSpace NameSpace `json:"namespace"`

	// Key is the unique name of the stored value.
	Key string `json:"key"`

	// Type indicates the data type of the value (string, float, boolean, binary).
	Type ValueType `json:"type"`

	// Size indicates the size of the stored value in bytes.
	Size int `json:"size"`
}

// ListValuesResponse represents the response returned by the Scan/List Values API.
type ListValuesResponse struct {
	// Count is the total number of items returned in this response.
	Count int `json:"count"`

	// Next is the offset value to be used for the next page of results.
	// If there are no more results, this may be 0 or equal to the total count.
	Next int `json:"next"`

	// Data contains the list of keys matching the scan criteria.
	Data []ValueItem `json:"data"`
}

// listValuesResponse is the internal struct used to decode the raw API JSON.
// It includes BaseResponse to handle success flags and error codes.
type listValuesResponse struct {
	common.BaseResponse
	Count int         `json:"count"`
	Next  int         `json:"next"`
	Data  []ValueItem `json:"data"`
}

// ScanAllAvailableItem retrieves a list of keys from the Value Store based on the provided filter and sorting options.
//
// It performs validation on the input scope and sort orders before making the API call.
//
// Parameters:
//   - ctx: The context for controlling the lifecycle of the request.
//   - input: A pointer to ListValuesRequest containing filter, sorting, and pagination options.
//
// Returns:
//   - *ListValuesResponse: The list of found keys and pagination details on success.
//   - error: An error if validation fails or the API call encounters an issue.
func (v *ValueStoreManagement) ScanAvailableItems(ctx context.Context, input *ListValuesRequest) (*ListValuesResponse, error) {

	// 1. Validate Input Struct
	if input == nil {
		return nil, &errors.AnedyaError{
			Message: "scan request cannot be nil",
			Err:     errors.ErrInputRequired,
		}
	}

	// Check if Namespace Scope is empty (Namespace is Required)
	if input.Filter.NameSpace.Scope == "" {
		return nil, &errors.AnedyaError{
			Message: "filter.namespace.scope is required",
			Err:     errors.ErrNamespaceScopeRequired,
		}
	}
	// Check if scope is valid
	if !isValidScope(input.Filter.NameSpace.Scope) {
		return nil, &errors.AnedyaError{
			Message: "invalid namespace scope",
			Err:     errors.ErrInvalidNamespaceScope,
		}
	}
	// Check if OrderBy is empty (OrderBy is Required)
	if input.OrderBy == "" {
		return nil, &errors.AnedyaError{
			Message: "orderby is required (valid values: namespace, key, created)",
			Err:     errors.ErrOrderByRequired,
		}
	}

	// Check if orderby is valid
	if !isValidOrderBy(input.OrderBy) {
		return nil, &errors.AnedyaError{
			Message: "Invalid order by",
			Err:     errors.ErrInvalidOrderBy,
		}
	}

	// Check if sort order is valid
	if input.Order != "" && !isValidSortOrder(input.Order) {
		return nil, &errors.AnedyaError{
			Message: "Invalid sort order",
			Err:     errors.ErrInvalidSortOrder,
		}
	}
	// 2. Encode Request
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to encode scan request",
			Err:     errors.ErrRequestEncodeFailed,
		}
	}

	// 3. Build HTTP Request
	url := fmt.Sprintf("%s/v1/valuestore/scan", v.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to build scan request",
			Err:     errors.ErrRequestBuildFailed,
		}
	}

	// 4. Execute Request
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to execute scan request",
			Err:     errors.ErrRequestFailed,
		}
	}
	defer resp.Body.Close()

	// 5. Read Response Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to read scan response",
			Err:     errors.ErrResponseReadFailed,
		}
	}

	// 6. Decode Response
	var apiResp listValuesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, &errors.AnedyaError{
			Message: "failed to decode scan response",
			Err:     errors.ErrResponseDecodeFailed,
		}
	}

	// 8. Handle API-level errors
	if !apiResp.Success {
		return nil, errors.GetError(apiResp.ReasonCode, apiResp.Error)
	}

	return &ListValuesResponse{
		Count: apiResp.Count,
		Next:  apiResp.Next,
		Data:  apiResp.Data,
	}, nil
}
