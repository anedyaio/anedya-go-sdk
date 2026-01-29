package accesstokens

import "net/http"

// AccessTokenManagement provides methods to create, revoke,
// and manage access tokens using the Anedya API.
//
// It wraps an HTTP client and a base URL used to construct requests.
type AccessTokenManagement struct {

	// httpClient is the underlying HTTP client used to
	// perform API requests.
	httpClient *http.Client

	// baseURL is the root API endpoint used for all
	// access token management requests.
	baseURL string
}

// NewAccessTokenManagement creates a new AccessTokenManagement client.
//
// The provided http.Client is used for all network communication,
// and baseURL specifies the API server address.
func NewAccessTokenManagement(c *http.Client, baseURL string) *AccessTokenManagement {
	return &AccessTokenManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}
