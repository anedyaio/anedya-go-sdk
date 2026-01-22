package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RemoveChildNodeRequest defines the request body for removing a specific child node
type RemoveChildNodeRequest struct {
	ParentId  string `json:"parentId"`  // Required: Parent node UUID
	ChildNode string `json:"childNode"` // Required: Child node UUID to remove
}

// RemoveChildNodeResponse defines the full API response (success + error)
type RemoveChildNodeResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode,omitempty"`
}

// RemoveChildNode removes a specific child node from a parent node
func (nm *NodeManagement) RemoveChildNode(ctx context.Context, req *RemoveChildNodeRequest) error {
	// Strong validation — fail fast with clear errors
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.ParentId == "" {
		return fmt.Errorf("parentId is required and cannot be empty")
	}
	if req.ChildNode == "" {
		return fmt.Errorf("childNode is required and cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/node/child/remove", nm.baseURL)

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
	var removeResp RemoveChildNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&removeResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error %d: %s (reason: %s)", resp.StatusCode, removeResp.Error, removeResp.ReasonCode)
	}

	// Application-level error
	if !removeResp.Success {
		return fmt.Errorf("failed to remove child node: %s (reason: %s)", removeResp.Error, removeResp.ReasonCode)
	}

	return nil
}
