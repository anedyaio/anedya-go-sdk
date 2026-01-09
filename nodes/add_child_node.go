package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChildNode represents a child node with optional alias
type ChildNode struct {
	NodeId string `json:"nodeId"`          // Required: Child node UUID
	Alias  string `json:"alias,omitempty"` // Optional: Human-readable name
}

// AddChildNodeRequest defines the request body for adding child nodes
type AddChildNodeRequest struct {
	ParentId   string      `json:"parentId"`   // Required: Parent node UUID
	ChildNodes []ChildNode `json:"childNodes"` // Required: At least one child
}

// AddChildNodeResponse defines the full API response (success + error)
type AddChildNodeResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode,omitempty"`
}

// AddChildNode adds one or more child nodes under a parent node
// The parent acts as a gateway/router for the children
func (nm *NodeManagement) AddChildNode(ctx context.Context, req *AddChildNodeRequest) error {
	// Strong validation — fail fast with clear, helpful errors
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.ParentId == "" {
		return fmt.Errorf("parentId is required and cannot be empty")
	}
	if len(req.ChildNodes) == 0 {
		return fmt.Errorf("childNodes array is required and must contain at least one child")
	}

	// Validate each child node
	for i, child := range req.ChildNodes {
		if child.NodeId == "" {
			return fmt.Errorf("childNodes[%d]: nodeId is required and cannot be empty", i)
		}
	}

	url := fmt.Sprintf("%s/v1/node/child/add", nm.baseURL)

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := nm.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Always decode full response first
	var addResp AddChildNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&addResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error %d: %s (reason: %s)", resp.StatusCode, addResp.Error, addResp.ReasonCode)
	}

	// Application-level error
	if !addResp.Success {
		return fmt.Errorf("failed to add child nodes: %s (reason: %s)", addResp.Error, addResp.ReasonCode)
	}

	return nil
}
