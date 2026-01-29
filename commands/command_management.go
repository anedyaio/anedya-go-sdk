// Package commands provides APIs to manage and execute
// commands on nodes within the Anedya platform.
package commands

import "net/http"

// CommandManagement provides methods to interact with
// command-related APIs.
//
// It encapsulates the HTTP client and base URL required
// to perform command execution and management operations.
type CommandManagement struct {
	httpClient *http.Client
	baseURL    string
}

// NewCommandManagement creates and returns a new instance
// of CommandManagement.
//
// Parameters:
//   - httpClient: Optional custom HTTP client. If nil,
//     http.DefaultClient is used.
//   - baseURL: Base URL of the Anedya API server.
//
// This function ensures a valid HTTP client is always
// available for making command-related API requests.
func NewCommandManagement(httpClient *http.Client, baseURL string) *CommandManagement {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &CommandManagement{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
