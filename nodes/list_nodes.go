package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//
// ===== Request & Response Models =====
//

// ListNodeRequest represents the request payload used
// to fetch a list of nodes from the platform.
//
// This struct supports pagination and ordering.
type ListNodeRequest struct {
	Limit  int    `json:"limit,omitempty"`  // Maximum number of nodes to return in one call
	Offset int    `json:"offset,omitempty"` // Offset for pagination (used to fetch next set)
	Order  string `json:"order,omitempty"`  // Sorting order: "asc" or "desc" (based on creation time)
}

// ListNodesResponse represents the complete API response
// returned by the List Nodes endpoint.
//
// Same structure is used for both success and error responses.
type ListNodesResponse struct {
	Success      bool     `json:"success"`                // Indicates whether API call was successful
	Error        string   `json:"error"`                  // Error message if success is false
	ReasonCode   string   `json:"reasonCode,omitempty"`   // Machine-readable error reason
	CurrentCount int      `json:"currentCount,omitempty"` // Number of nodes returned in this response
	TotalCount   int      `json:"totalCount,omitempty"`   // Total number of nodes available
	Nodes        []string `json:"nodes,omitempty"`        // List of node IDs
	Offset       int      `json:"offset,omitempty"`       // Offset value used for this response
}

//
// ===== Public SDK Method =====
//

// ListNodes fetches a paginated list of node IDs from the Anedya platform.
//
// Nodes are returned in the order of creation time.
// By default, newest nodes are returned first unless order is specified.
//
// Parameters:
//   - ctx: Context for request cancellation, timeout, and tracing.
//   - req: ListNodeRequest containing pagination and ordering details.
//
// Returns:
//   - *ListNodesResponse: Contains list of node IDs and pagination metadata.
//   - error: Non-nil error if request fails, API returns error, or response is invalid.
//
// Example usage:
//
//	resp, err := nodeManager.ListNodes(ctx, &ListNodeRequest{Limit: 100, Offset: 0})
//	if err != nil {
//	    log.Fatal(err)
//	}
func (nm *NodeManagement) ListNodes(
	ctx context.Context,
	req *ListNodeRequest,
) (*ListNodesResponse, error) {

	// Set default limit if user has not provided one
	// This avoids accidental large data fetches
	if req.Limit == 0 {
		req.Limit = 1000 // Safe default value
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/list", nm.baseURL)

	// Convert request payload to JSON
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP POST request with context
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Execute HTTP request
	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode API response (works for both success and error cases)
	var listResp ListNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Handle non-200 HTTP responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"API error %d: %s (reason: %s)",
			resp.StatusCode,
			listResp.Error,
			listResp.ReasonCode,
		)
	}

	// Handle logical failure returned by API
	if !listResp.Success {
		return nil, fmt.Errorf(
			"failed to fetch node list: %s (reason: %s)",
			listResp.Error,
			listResp.ReasonCode,
		)
	}

	// Return successful response
	return &listResp, nil
}
