package client

import (
	"net/http"
	"time"

	"github.com/anedyaio/anedya-go-sdk/dataAccess"
	"github.com/anedyaio/anedya-go-sdk/nodes"
	"github.com/anedyaio/anedya-go-sdk/variable"
)

type Client struct {
	NodeManagement     *nodes.NodeManagement
	VariableManagement *variable.VariableManagement
	DataManagement     *dataAccess.DataManagement
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
		NodeManagement:     nodes.NewNodeManagement(hc, baseURL),
		VariableManagement: variable.NewVariableManagement(hc, baseURL),
		DataManagement:     dataAccess.NewDataManagement(hc, baseURL),
	}
}
