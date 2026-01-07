package node

import "net/http"

type Node struct {
	httpClient *http.Client
	baseURL    string
}

func GetNode(httpClient *http.Client, baseURL string) *Node {
	return &Node{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
