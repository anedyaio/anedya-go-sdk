// Package aggregations provides a client and configuration
// types for interacting with aggregation-related APIs
// exposed by the Anedya platform.
package aggregations

import "net/http"

// AggregationsClient provides methods to interact with
// aggregation APIs on the Anedya platform.
type AggregationsClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewAggregationsClient creates and returns a new AggregationsClient.
//
// Parameters:
//   - c: Optional custom HTTP client. If nil, http.DefaultClient is used.
//   - baseURL: Base URL of the Anedya API
//     (for example, "https://api.ap-in-1.anedya.io").
//
// Returns:
//   - *AggregationsClient: A configured client ready to make
//     aggregation API calls.
func NewAggregationsClient(c *http.Client, baseURL string) *AggregationsClient {
	if c == nil {
		c = http.DefaultClient
	}

	return &AggregationsClient{
		httpClient: c,
		baseURL:    baseURL,
	}
}
