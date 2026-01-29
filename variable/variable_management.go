package variable

import "net/http"

// VariableManagement provides methods to create, update, delete,
// and retrieve variables from the Anedya API.
//
// It wraps an HTTP client and a base URL used to construct requests.
type VariableManagement struct {

	// httpClient is the underlying HTTP client used to
	// perform API requests.
	httpClient *http.Client

	// baseURL is the root API endpoint used for all
	// variable management requests.
	baseURL string
}

// NewVariableManagement creates a new VariableManagement client.
//
// The provided http.Client is used for all network communication,
// and baseURL specifies the API server address.
func NewVariableManagement(c *http.Client, baseURL string) *VariableManagement {
	return &VariableManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}
