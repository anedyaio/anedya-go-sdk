package client

import (
	"net/http"
	"time"

	accesstokens "github.com/anedyaio/anedya-go-sdk/accessTokens"
	"github.com/anedyaio/anedya-go-sdk/nodes"
	"github.com/anedyaio/anedya-go-sdk/variable"
)

type Client struct {
	NodeManagement        *nodes.NodeManagement
	VariableManagement    *variable.VariableManagement
	AccessTokenManagement *accesstokens.AccessTokenManagement
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
		NodeManagement:        nodes.NewNodeManagement(hc, baseURL),
		VariableManagement:    variable.NewVariableManagement(hc, baseURL),
		AccessTokenManagement: accesstokens.NewAccessTokenManagement(hc, baseURL),
	}
}
