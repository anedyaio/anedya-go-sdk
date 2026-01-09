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

// getNodeDetailsRequest represents the request body sent to
// the Node Details API.
//
// It contains a list of node IDs for which details are required.
type getNodeDetailsRequest struct {
	Nodes []string `json:"nodes"`
}

// nodeDetailsDTO represents the node details object
// received from the API response.
//
// This is a Data Transfer Object (DTO) and is later
// converted into the SDK's Node model.
type nodeDetailsDTO struct {
	NodeId          string `json:"nodeId"`     // Unique Node ID
	NodeName        string `json:"node_name"`  // Human readable node name
	NodeDescription string `json:"node_desc"`  // Description of the node
	Tags            []Tag  `json:"tags"`       // Associated tags
	PreauthId       string `json:"preauth_id"` // Pre-authorisation ID
}

// getNodeDetailsResponse represents the complete response
// returned by the Node Details API.
type getNodeDetailsResponse struct {
	Success    bool                      `json:"success"`              // Indicates request success
	Error      string                    `json:"error"`                // Error message (if any)
	ReasonCode string                    `json:"reasonCode,omitempty"` // Error reason code
	Data       map[string]nodeDetailsDTO `json:"data,omitempty"`       // Node data mapped by node ID
}

//
// ===== Public SDK Method =====
//

// GetNodeDetails fetches details of a single node from the server.
//
// Parameters:
//   - ctx: Context used for request cancellation, timeout, and tracing.
//   - nodeID: Unique identifier of the node whose details are required.
//
// Returns:
//   - *Node: Fully populated Node object if the API call is successful.
//   - error: Non-nil error if validation fails, API fails, or node is not found.
//
// Example usage:
//
//	node, err := nodeManager.GetNodeDetails(ctx, "node-123")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (nm *NodeManagement) GetNodeDetails(
	ctx context.Context,
	nodeID string,
) (*Node, error) {

	// Validate input parameter
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}

	// Construct API endpoint URL
	url := fmt.Sprintf("%s/v1/node/details", nm.baseURL)

	// Prepare request payload
	reqBody := getNodeDetailsRequest{
		Nodes: []string{nodeID},
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(reqBody)
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

	// Decode API response
	var apiResp getNodeDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Handle non-200 HTTP responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"API error %d: %s (reason: %s)",
			resp.StatusCode,
			apiResp.Error,
			apiResp.ReasonCode,
		)
	}

	// Handle logical API failure
	if !apiResp.Success {
		return nil, fmt.Errorf(
			"failed to get node details: %s (reason: %s)",
			apiResp.Error,
			apiResp.ReasonCode,
		)
	}

	// Extract node details from response map
	dto, ok := apiResp.Data[nodeID]
	if !ok {
		return nil, fmt.Errorf("node %s not found in response", nodeID)
	}

	// Convert DTO to SDK Node model and return
	return &Node{
		nodeManagement:  nm,
		NodeId:          dto.NodeId,
		NodeName:        dto.NodeName,
		NodeDescription: dto.NodeDescription,
		Tags:            dto.Tags,
		PreauthId:       dto.PreauthId,
	}, nil
}
