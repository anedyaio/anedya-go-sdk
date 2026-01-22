package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DeleteNodeRequest defines the request body for deleting a node
type DeleteNodeRequest struct {
	NodeID string `json:"nodeid"` // Required: Node UUID to delete
}

// DeleteNodeResponse defines the full API response (used for both success and error)
type DeleteNodeResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode,omitempty"`
}

// DeleteNode permanently deletes a node and ALL its associated data
// Includes telemetry, configs, child nodes, submissions — everything!
// THIS IS IRREVERSIBLE. Use only when absolutely certain.
func (nm *NodeManagement) DeleteNode(ctx context.Context, req *DeleteNodeRequest) error {
	// Strong validation — prevent empty or invalid calls early
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.NodeID == "" {
		return fmt.Errorf("nodeid is required and cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/node/delete", nm.baseURL)

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
	var deleteResp DeleteNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// HTTP-level error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error %d: %s (reason: %s)", resp.StatusCode, deleteResp.Error, deleteResp.ReasonCode)
	}

	// Application-level error
	if !deleteResp.Success {
		return fmt.Errorf("node deletion failed: %s (reason: %s)", deleteResp.Error, deleteResp.ReasonCode)
	}

	return nil
}
