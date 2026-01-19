package nodes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Node represents a device or logical entity registered on the Anedya platform.
//
// A Node contains all basic information such as name, identifiers, binding status, and tags.
// It also holds a reference to NodeManagement, allowing it to perform node-level operations
// like fetching details, updating, or managing child nodes directly.
type Node struct {
	NodeId          string `json:"nodeId,omitempty"`          // Unique identifier for this node
	NodeName        string `json:"nodeName,omitempty"`        // Human-readable name for the node
	NodeDescription string `json:"nodeDescription,omitempty"` // Optional description of the node
	NodeIdentifier  string `json:"nodeIdentifier,omitempty"`  // Unique external identifier, if any
	BindingStatus   bool   `json:"bindingStatus,omitempty"`   // Whether the node is bound to a device
	NodeBindingKey  string `json:"nodeBindingKey,omitempty"`  // Key used for binding operations
	ConnectionKey   string `json:"connectionKey,omitempty"`   // Key used for device connections
	CreatedAt       string `json:"createdAt,omitempty"`       // Creation timestamp as string
	Suspended       bool   `json:"suspended,omitempty"`       // Node suspension status
	Modified        string `json:"modified,omitempty"`        // Last modification timestamp
	Tags            []Tag  `json:"tags,omitempty"`            // Optional list of tags for categorisation
	PreauthId       string `json:"preauthId,omitempty"`       // Preauthorization ID for node

	// nodeManagement is an internal reference to the NodeManagement client.
	// It is required for all node-related API calls.
	nodeManagement *NodeManagement `json:"-"`
}

// Tag represents a key-value metadata pair attached to a node.
//
// Tags are used for categorization, filtering, or adding extra metadata to nodes.
type Tag struct {
	Key   string `json:"key"`   // Tag name
	Value string `json:"value"` // Tag value
}

// NodeManagement manages all API interactions related to nodes.
//
// It stores the HTTP client and base URL required to communicate with
// the Anedya backend for node-related operations.
type NodeManagement struct {
	httpClient *http.Client // HTTP client used for API requests
	baseURL    string       // Base URL for node endpoints
}

// NewNodeManagement creates a new NodeManagement instance.
//
// Parameters:
//   - c: HTTP client to be reused for all node operations
//   - baseURL: Base API URL for node-related endpoints
//
// Returns:
//   - *NodeManagement: initialized NodeManagement instance
func NewNodeManagement(c *http.Client, baseURL string) *NodeManagement {
	return &NodeManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// ==================== Node Wrapper Methods ====================

// GetDetails fetches the latest details of the node from the server.
//
// This method performs the following:
//  1. Validates that the NodeManagement client is initialized.
//  2. Calls the NodeManagement's GetNodeDetails API with the NodeId.
//  3. Updates the current Node instance with the latest details.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//
// Returns:
//   - *Node: Updated node object with fresh data
//   - error: Error if NodeManagement is nil, node not found, or API call fails
func (n *Node) GetDetails(ctx context.Context) (*Node, error) {
	if n.nodeManagement == nil {
		return nil, errors.New("node management client is not initialized")
	}

	req := &GetNodeDetailsRequest{
		Nodes: []string{n.NodeId},
	}

	data, err := n.nodeManagement.GetNodeDetails(ctx, req)
	if err != nil {
		return nil, err
	}

	details, ok := data[n.NodeId]
	if !ok {
		return nil, errors.New("node details not found")
	}

	n.NodeName = details.NodeName
	n.NodeDescription = details.NodeDescription
	n.NodeIdentifier = details.NodeIdentifier
	n.BindingStatus = details.BindingStatus
	n.NodeBindingKey = details.NodeBindingKey
	n.ConnectionKey = details.ConnectionKey
	n.CreatedAt = details.CreatedAt
	n.Suspended = details.Suspended
	n.Modified = details.Modified
	n.Tags = details.Tags

	return n, nil
}

// ListChildNodes returns child nodes attached to this node.
//
// This method performs the following:
//  1. Validates that NodeManagement client is initialized.
//  2. Calls the ListChildNodes API with parent NodeId, limit, and offset for pagination.
//  3. Returns a slice of child Node instances with the same NodeManagement reference.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - limit: Maximum number of child nodes to return
//   - offset: Pagination offset
//
// Returns:
//   - []*Node: Slice of child nodes
//   - error: Error if NodeManagement is nil or API call fails
func (n *Node) ListChildNodes(ctx context.Context, limit int, offset int) ([]*Node, error) {
	if n.nodeManagement == nil {
		return nil, errors.New("node management client is not initialized")
	}

	req := &ListChildNodesRequest{
		ParentId: n.NodeId,
		Limit:    limit,
		Offset:   offset,
	}

	resp, err := n.nodeManagement.ListChildNodes(ctx, req)
	if err != nil {
		return nil, err
	}

	nodes := make([]*Node, 0, len(resp.Data))
	for _, child := range resp.Data {
		node := &Node{
			NodeId:         child.ChildId,
			NodeName:       child.Alias,
			CreatedAt:      fmt.Sprintf("%d", child.CreatedAt),
			nodeManagement: n.nodeManagement,
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateNode applies updates to the node.
//
// This method performs the following:
//  1. Validates that NodeManagement client is initialized.
//  2. Calls NodeManagement's UpdateNode API with the provided updates.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - updates: Slice of NodeUpdate containing fields to update
//
// Returns:
//   - error: Error if NodeManagement is nil or API call fails
func (n *Node) UpdateNode(ctx context.Context, updates []NodeUpdate) error {
	if n.nodeManagement == nil {
		return errors.New("node management client is not initialized")
	}

	req := &UpdateNodeRequest{
		NodeID:  n.NodeId,
		Updates: updates,
	}

	return n.nodeManagement.UpdateNode(ctx, req)
}

// AuthorizeDevice authorizes a device to connect with this node.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - deviceID: Unique device identifier to authorize
//
// Returns:
//   - error: Error if NodeManagement is nil, deviceID is empty, or API call fails
func (n *Node) AuthorizeDevice(ctx context.Context, deviceID string) error {
	if n.nodeManagement == nil {
		return errors.New("node management client is not initialized")
	}

	if deviceID == "" {
		return errors.New("deviceID is required")
	}

	req := &AuthorizeDeviceRequest{
		NodeID:   n.NodeId,
		DeviceID: deviceID,
	}

	return n.nodeManagement.AuthorizeDevice(ctx, req)
}

// AddChildNode attaches one or more child nodes to this node.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - childNodes: Slice of ChildNodeRequest containing child IDs and aliases
//
// Returns:
//   - error: Error if NodeManagement is nil or API call fails
func (n *Node) AddChildNode(ctx context.Context, childNodes []ChildNodeRequest) error {
	if n.nodeManagement == nil {
		return errors.New("node management client is not initialized")
	}

	req := &AddChildNodeRequest{
		ParentId:   n.NodeId,
		ChildNodes: childNodes,
	}

	return n.nodeManagement.AddChildNode(ctx, req)
}

// ClearChildNodes removes all child nodes attached to this node.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//
// Returns:
//   - error: Error if NodeManagement is nil or API call fails
func (n *Node) ClearChildNodes(ctx context.Context) error {
	if n.nodeManagement == nil {
		return errors.New("node management client is not initialized")
	}

	req := &ClearChildNodesRequest{
		ParentId: n.NodeId,
	}

	return n.nodeManagement.ClearChildNodes(ctx, req)
}

// GetConnectionKey retrieves the connection key for this node.
//
// This key is generally used by devices to establish connections.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//
// Returns:
//   - string: Connection key
//   - error: Error if NodeManagement is nil or API call fails
func (n *Node) GetConnectionKey(ctx context.Context) (string, error) {
	if n.nodeManagement == nil {
		return "", errors.New("node management client is not initialized")
	}

	req := &GetConnectionKeyRequest{
		NodeID: n.NodeId,
	}

	return n.nodeManagement.GetConnectionKey(ctx, req)
}

// RemoveChildNode detaches a specific child node from this node.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - childNodeID: NodeId of the child node to remove
//
// Returns:
//   - error: Error if NodeManagement is nil or API call fails
func (n *Node) RemoveChildNode(ctx context.Context, childNodeID string) error {
	if n.nodeManagement == nil {
		return errors.New("node management client is not initialized")
	}

	req := &RemoveChildNodeRequest{
		ParentId:  n.NodeId,
		ChildNode: childNodeID,
	}

	return n.nodeManagement.RemoveChildNode(ctx, req)
}
