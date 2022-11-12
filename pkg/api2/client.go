// Massdriver GraphQL API queries/mutations using Genqlient
//
//go:generate go run github.com/Khan/genqlient
package api2

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

func NewClient(apiKey string) graphql.Client {
	c := http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport, apiKey: apiKey}}
	return graphql.NewClient("https://api.massdriver.cloud/api/", &c)
}

type authedTransport struct {
	wrapped http.RoundTripper
	apiKey  string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("authorization", "Bearer "+t.apiKey)
	return t.wrapped.RoundTrip(req)
}
