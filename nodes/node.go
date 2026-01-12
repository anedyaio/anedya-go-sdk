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
