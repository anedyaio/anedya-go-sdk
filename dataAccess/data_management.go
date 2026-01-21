package dataAccess

import "net/http"

// DataManagement provides methods to interact with
// data-related APIs exposed by the platform.
//
// It encapsulates the HTTP client and base URL required
// to perform all data management operations.
type DataManagement struct {
	httpClient *http.Client
	baseURL    string
}

// NewDataManagement creates and returns a new instance of DataManagement.
//
// Parameters:
//   - httpClient: Optional custom HTTP client. If nil, http.DefaultClient is used.
//   - baseURL: Base URL of the API server.
//
// This function ensures a valid HTTP client is always available
// for making API requests.
func NewDataManagement(httpClient *http.Client, baseURL string) *DataManagement {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &DataManagement{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
