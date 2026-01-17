package nodes

import (
	"net/http"
)

// Node represents a device or logical entity in the Anedya platform.
//
// The Node struct contains metadata associated with a node and
// maintains a reference to NodeManagement for performing
// node-specific operations.
type Node struct {

	// NodeId is the unique identifier assigned to the node
	// by the Anedya platform.
	NodeId string `json:"nodeId,omitempty"`

	// NodeName is the human-readable name of the node.
	NodeName string `json:"nodeName,omitempty"`

	// NodeDescription is an optional description
	// providing additional information about the node.
	NodeDescription string `json:"nodeDescription,omitempty"`

	// Tags is an optional list of metadata key-value pairs
	// associated with the node.
	Tags []Tag `json:"tags,omitempty"`

	// PreauthId is an optional pre-authorization identifier
	// associated with the node.
	PreauthId string `json:"preauthId,omitempty"`

	// nodeManagement is the internal client used to perform
	// operations related to this node.
	//
	// This field is intentionally excluded from JSON serialization.
	nodeManagement *NodeManagement `json:"-"`
}

// Tag represents a key-value metadata pair
// used to categorize or describe a node.
type Tag struct {

	// Key is the tag identifier.
	Key string `json:"key"`

	// Value is the tag value.
	Value string `json:"value"`
}

// NodeManagement provides methods to create, update,
// delete, and manage nodes in the Anedya platform.
//
// It encapsulates the HTTP client and base API URL
// required for node-related API operations.
type NodeManagement struct {

	// httpClient is used to make HTTP requests
	// to the Anedya API.
	httpClient *http.Client

	// baseURL is the root endpoint for
	// all node-related API calls.
	baseURL string
}

// NewNodeManagement creates and returns a new
// NodeManagement client instance.
//
// The provided http.Client is used for all network requests,
// and baseURL specifies the root API endpoint for node operations.
func NewNodeManagement(c *http.Client, baseURL string) *NodeManagement {
	return &NodeManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}
