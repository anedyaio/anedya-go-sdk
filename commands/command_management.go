package commands

import "net/http"

// CommandManagement provides APIs for sending, listing,
// invalidating, and querying command status for nodes/devices
// in the Anedya platform.
type CommandManagement struct {
	// httpClient is used to execute HTTP requests.
	httpClient *http.Client

	// baseURL is the root endpoint for all command APIs.
	baseURL string
}

// NewCommandManagement creates and returns a new CommandManagement client.
//
// Parameters:
//   - httpClient: Optional custom HTTP client. If nil, http.DefaultClient is used.
//   - baseURL: Base URL of the Anedya API server.
//
// This constructor ensures a valid HTTP client is always set
// and allows consumers to inject custom configurations such as
// timeouts, transports, or middleware.
func NewCommandManagement(httpClient *http.Client, baseURL string) *CommandManagement {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &CommandManagement{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
