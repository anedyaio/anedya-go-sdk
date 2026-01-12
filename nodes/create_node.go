package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateNodeRequest defines the request body for creating a node
type CreateNodeRequest struct {
	NodeName  string `json:"node_name"`
	NodeDesc  string `json:"node_desc,omitempty"`
	Tags      []Tag  `json:"tags,omitempty"`
	PreauthId string `json:"preauth_id,omitempty"`
}

// CreateNodeResponse defines the API response (used for both success and error)
type CreateNodeResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	ReasonCode string `json:"reasonCode,omitempty"`
	NodeId     string `json:"nodeId,omitempty"`
}

// CreateNode creates a new node on the Anedya platform
func (nm *NodeManagement) CreateNode(ctx context.Context, req *CreateNodeRequest) (*Node, error) {
	url := fmt.Sprintf("%s/v1/node/create", nm.baseURL)

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

	// Always decode the full response — works for success and error cases
	var createResp CreateNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// First check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s (reason: %s)", resp.StatusCode, createResp.Error, createResp.ReasonCode)
	}

	// Then check application-level success
	if !createResp.Success {
		return nil, fmt.Errorf("node creation failed: %s (reason: %s)", createResp.Error, createResp.ReasonCode)
	}

	// Success: return populated Node
	return &Node{
		nodeManagement:  nm,
		NodeId:          createResp.NodeId,
		NodeName:        req.NodeName,
		NodeDescription: req.NodeDesc,
		Tags:            req.Tags,
		PreauthId:       req.PreauthId,
	}, nil
}
