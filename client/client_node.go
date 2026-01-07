package client

import (
	"context"

	"github.com/anedyaio/anedya-go-sdk/node"
)

func (c *Client) CreateNode(ctx context.Context, nodeName, nodeDesc string, tags map[string]string) (*node.Node, error) {

	node := node.GetNode(c.httpClient, c.baseURL)

	// Implement the logic to create a node using REST API calls here.

	return node, nil
}
