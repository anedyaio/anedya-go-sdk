package health

import "net/http"

// HealthManagement provides methods to interact with
// health-related APIs on the Anedya platform.
//
// It is responsible for checking node connectivity,
// heartbeat status, and overall device health.
type HealthManagement struct {
	// httpClient is the HTTP client used to make API requests.
	httpClient *http.Client

	// baseURL is the base endpoint for all health APIs.
	baseURL string
}

// NewHealthManagement creates and returns a new HealthManagement client.
//
// Parameters:
//   - c: Custom HTTP client to be used for API requests.
//     If nil, http.DefaultClient should be provided by the caller.
//   - baseURL: Base URL of the Anedya API server.
//
// Returns:
//   - *HealthManagement: A configured health API client instance.
func NewHealthManagement(c *http.Client, baseURL string) *HealthManagement {
	return &HealthManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}
