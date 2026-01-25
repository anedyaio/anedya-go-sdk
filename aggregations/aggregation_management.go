package aggregations

import "net/http"

// AggregationsManagement provides methods to interact with
// aggregation-related APIs on the Anedya platform.
//
// It is used to fetch and compute aggregated data
// across nodes, variables, and time ranges.
type AggregationsManagement struct {
	// httpClient is the HTTP client used to make API requests.
	httpClient *http.Client

	// baseURL is the base endpoint for all aggregation APIs.
	baseURL string
}

// NewAggregationsManagement creates and returns a new
// AggregationsManagement client.
//
// Parameters:
//   - c: Custom HTTP client used for API requests.
//     If nil, http.DefaultClient should be used by the caller.
//   - baseURL: Base URL of the Anedya API server
//     (e.g., "https://api.ap-in-1.anedya.io").
//
// Returns:
//   - *AggregationsManagement: A configured aggregation API client.
func NewAggregationsManagement(c *http.Client, baseURL string) *AggregationsManagement {
	if c == nil {
		c = http.DefaultClient
	}
	return &AggregationsManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}
