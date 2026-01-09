package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListChildNodesRequest defines the request body for listing child nodes
type ListChildNodesRequest struct {
	ParentId string `json:"parentId"`         // Required: Parent node UUID
	Limit    int    `json:"limit,omitempty"`  // Optional: Max items (default usually 100)
	Offset   int    `json:"offset,omitempty"` // Optional: Pagination offset
}

// ListChildNodesResponse defines the full API response (success + error)
type ListChildNodesResponse struct {
	Success    bool    `json:"success"`
	Error      string  `json:"error"`
	ReasonCode string  `json:"reasonCode,omitempty"`
	TotalCount int     `json:"totalCount,omitempty"`
	Count      int     `json:"count,omitempty"`
	Next       int     `json:"next,omitempty"`
	Data       []*Node `json:"data,omitempty"`
}

// ListAllChildNodes fetches a paginated list of child nodes for a parent node
func (nm *NodeManagement) ListAllChildNodes(ctx context.Context, req *ListChildNodesRequest) (*ListChildNodesResponse, error) {
	// Strong validation — fail fast if parentId missing
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.ParentId == "" {
		return nil, fmt.Errorf("parentId is required and cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/node/child/list", nm.baseURL)

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Always decode full response first
	var listResp ListChildNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s (reason: %s)", resp.StatusCode, listResp.Error, listResp.ReasonCode)
	}

	// Application-level error
	if !listResp.Success {
		return nil, fmt.Errorf("failed to list child nodes: %s (reason: %s)", listResp.Error, listResp.ReasonCode)
	}

	// Link nodeManagement to each child Node for direct method calls
	for _, child := range listResp.Data {
		if child != nil {
			child.nodeManagement = nm
		}
	}

	return &listResp, nil
}
