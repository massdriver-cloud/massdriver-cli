package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method string
	Path   string
	Body   io.Reader
}

func NewRequest(method string, path string, body io.Reader) *Request {
	req := new(Request)

	req.Method = method
	req.Path = path
	req.Body = body

	return req
}

func (req *Request) toHTTPRequest(c *MassdriverClient) (*http.Request, error) {
	url, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	url.Path = req.Path

	httpReq, err := http.NewRequest(req.Method, url.String(), req.Body)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		httpReq.Header.Set("X-Md-Api-Key", c.apiKey)
	} else {
		fmt.Println("Warning: API Key not specified")
	}
	// for now assuming everything is json
	httpReq.Header.Set("Content-Type", "application/json")

	return httpReq, nil
}
