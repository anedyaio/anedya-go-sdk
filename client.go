package client

import (
	"net/http"
	"time"

	"github.com/anedyaio/anedya-go-sdk/nodes"
)

type Client struct {
	NodeManagement *nodes.NodeManagement
}

func NewClient(apiKey, baseURL string) *Client {

	auth := &authTransport{
		apiKey: apiKey,
		next:   http.DefaultTransport,
	}

	hc := &http.Client{
		Timeout:   30 * time.Second,
		Transport: auth,
	}

	return &Client{
		NodeManagement: nodes.NewNodeManagement(hc, baseURL),
	}
}
