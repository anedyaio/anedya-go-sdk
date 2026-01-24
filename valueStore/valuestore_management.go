package valuestore

import "net/http"

// ValueStoreManagement provides methods to interact with the
// Anedya Value Store APIs.
//
// The client wraps an HTTP client and a base API URL which are used
// internally to construct and execute all API requests.
type ValueStoreManagement struct {

	// httpClient is the underlying HTTP client used to perform
	// network requests to the Anedya API.
	httpClient *http.Client

	// baseURL is the root API endpoint used for all
	// value store operations.
	baseURL string
}

// NewValueStoreMangement creates and returns a new ValueStoreManagement client.
//
// The provided http.Client is used for all network communication,
// and baseURL specifies the API server address.
func NewValueStoreManagement(c *http.Client, baseURL string) *ValueStoreManagement {
	return &ValueStoreManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}
