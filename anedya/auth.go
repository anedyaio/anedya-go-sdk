package anedya

import (
	"fmt"
	"net/http"
)

type authTransport struct {
	apiKey string
	next   http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newReq := req.Clone(req.Context())
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))
	newReq.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return t.next.RoundTrip(newReq)
}
