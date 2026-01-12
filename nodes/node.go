package nodes

import (
	"net/http"
)

// Node reperesents a node in Anedya Platform
type Node struct {
	nodeManagement *NodeManagement

	NodeId          string `json:"nodeId,omitempty"`
	NodeName        string `json:"nodeName,omitempty"`
	NodeDescription string `json:"nodeDescription,omitempty"`
	Tags            []Tag  `json:"tags,omitempty"`
	PreauthId       string `json:"preauthId,omitempty"`
}

// Tag represents key-value tag for node
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NodeManagement struct {
	httpClient *http.Client
	baseURL    string
}

func NewNodeManagement(c *http.Client, baseURL string) *NodeManagement {
	return &NodeManagement{
		httpClient: c,
		baseURL:    baseURL,
	}
}

// func (nm *NodeManagement) CreateNode(ctx context.Context, node *Node) (*Node, error) {

// 	// implement the logic for http request

// 	return &Node{
// 		nodeManagement: nm,
// 	}, nil
// }

// func (nm *NodeManagement) GetNode(ctx context.Context, nodeId string) (*Node, error) {
// 	return &Node{
// 		nodeManagement: nm,
// 	}, nil
// }

// func (nm *NodeManagement) UpdateNodeName(ctx context.Context, nodeName string) (*Node, error) {
// 	return &Node{
// 		nodeManagement: nm,
// 	}, nil
// }

// func (nm *NodeManagement) UpdateNodeName(ctx context.Context, nodeName string) (*Node, error) {
// 	return &Node{
// 		nodeManagement: nm,
// 	}, nil
// }

// ===========================
// Node Level Methods
// ===========================
// func (n *Node) UpdateNodeName(ctx context.Context, name string) error {
// 	n.nodeManagement.UpdateNodeName(ctx, name)
// 	return nil
// }
