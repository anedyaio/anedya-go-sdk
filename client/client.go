package client

import "net/http"

type Client struct {
	// Add client fields here
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL, apiKey string) *Client {

	auth := &authTransport{
		apiKey: apiKey,
		next:   http.DefaultTransport,
	}

	httpClient := &http.Client{
		Transport: auth,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
