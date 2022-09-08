package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

func NewClient() *graphql.Client {
	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}}
	client := graphql.NewClient("https://api.massdriver.cloud/api/", &c)
	return client
}

func NewSubscriptionClient() *graphql.SubscriptionClient {
	client := graphql.NewSubscriptionClient("https://api.massdriver.cloud/api/")
	return client
}

type transport struct {
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	apiKey := os.Getenv("MASSDRIVER_API_KEY")

	if apiKey == "" {
		log.Fatal().Msg("MASSDRIVER_API_KEY must be set")
	}

	bearer := fmt.Sprintf("Bearer %s", apiKey)
	req.Header.Add("authorization", bearer)
	return t.underlyingTransport.RoundTrip(req)
}
