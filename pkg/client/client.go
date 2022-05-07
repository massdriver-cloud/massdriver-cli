package client

import (
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MassdriverClient struct {
	client   HTTPClient
	endpoint string
	apiKey   string
}

var (
	MASSDRIVER_ENDPOINT = "https://api.massdriver.cloud"
)

func NewClient() *MassdriverClient {
	c := new(MassdriverClient)

	c.client = http.DefaultClient
	c.endpoint = MASSDRIVER_ENDPOINT
	c.apiKey = getApiKey()

	return c
}

// eventually this could walk through multiple sources (environment, then config file, etc)
func getApiKey() string {
	return os.Getenv("MASSDRIVER_API_KEY")
}

func (c *MassdriverClient) WithApiKey(apiKey string) *MassdriverClient {
	c.apiKey = apiKey
	return c
}

func (c *MassdriverClient) WithEndpoint(endpoint string) *MassdriverClient {
	c.endpoint = endpoint
	return c
}

func (c *MassdriverClient) Do(req *Request) (*http.Response, error) {
	httpReq, err := req.toHTTPRequest(c)
	if err != nil {
		return nil, err
	}

	return c.client.Do(httpReq)
}
