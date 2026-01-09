package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UpdateType represents valid update operations
type UpdateType string

const (
	UpdateNodeName UpdateType = "node_name"
	UpdateNodeDesc UpdateType = "node_desc"
	RegenerateKey  UpdateType = "key"       // Regenerate connection key
	DeleteTag      UpdateType = "deletetag" // Delete a tag by key
)

// UpdateItem represents one update operation
type UpdateItem struct {
	Type  UpdateType `json:"type"`
	Value string     `json:"value,omitempty"` // Required for node_name, node_desc
	Tag   *Tag       `json:"tag,omitempty"`   // Required for deletetag
}

// UpdateNodeRequest defines the full request body
type UpdateNodeRequest struct {
	NodeID  string       `json:"nodeid"`
	Updates []UpdateItem `json:"updates"`
}

// UpdateNodeResponse defines the full API response (success + error)
type UpdateNodeResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode,omitempty"`
}

// UpdateNodeDetails applies one or more updates to a node in a single call
func (nm *NodeManagement) UpdateNodeDetails(ctx context.Context, req *UpdateNodeRequest) error {
	// Strong validation — fail early with clear errors
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.NodeID == "" {
		return fmt.Errorf("nodeid is required and cannot be empty")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("at least one update must be provided in 'updates' array")
	}

	// Optional: extra validation for each update item
	for i, item := range req.Updates {
		switch item.Type {
		case UpdateNodeName, UpdateNodeDesc:
			if item.Value == "" {
				return fmt.Errorf("update item %d: 'value' is required for type '%s'", i, item.Type)
			}
		case DeleteTag:
			if item.Tag == nil || item.Tag.Key == "" {
				return fmt.Errorf("update item %d: 'tag.key' is required for type 'deletetag'", i)
			}
		case RegenerateKey:
			// No extra fields needed
		default:
			return fmt.Errorf("update item %d: invalid type '%s'", i, item.Type)
		}
	}

	url := fmt.Sprintf("%s/v1/node/update", nm.baseURL)

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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Some APIs return empty body on success (especially for key regenerate)
	if len(bodyBytes) == 0 {
		if resp.StatusCode == http.StatusOK {
			return nil // Success with no content
		}
		return fmt.Errorf("empty response body with non-200 status: %d", resp.StatusCode)
	}

	// Try to parse JSON response
	var updateResp UpdateNodeResponse
	if err := json.Unmarshal(bodyBytes, &updateResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w (raw body: %s)", err, string(bodyBytes))
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error %d: %s (reason: %s)", resp.StatusCode, updateResp.Error, updateResp.ReasonCode)
	}

	// Application-level error
	if !updateResp.Success {
		return fmt.Errorf("node update failed: %s (reason: %s)", updateResp.Error, updateResp.ReasonCode)
	}

	return nil
}
