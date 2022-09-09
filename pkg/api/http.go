package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
)

// TODO remove this once SA token working
// for now this can be populated with gql`query me { me { token } }`
var removeMeToken = "SFMyNTY.g2gDdAAAAA5kAAhfX21ldGFfX3QAAAAGZAAKX19zdHJ1Y3RfX2QAG0VsaXhpci5FY3RvLlNjaGVtYS5NZXRhZGF0YWQAB2NvbnRleHRkAANuaWxkAAZwcmVmaXhkAANuaWxkAAZzY2hlbWFkACNFbGl4aXIuTWFzc2RyaXZlci5NZWF0c3BhY2UuQWNjb3VudGQABnNvdXJjZW0AAAAIYWNjb3VudHNkAAVzdGF0ZWQABWJ1aWx0ZAAKX19zdHJ1Y3RfX2QAI0VsaXhpci5NYXNzZHJpdmVyLk1lYXRzcGFjZS5BY2NvdW50ZAALYXR0cmlidXRpb25kAANuaWxkABNiZXRhX2FjY2Vzc19lbmFibGVkZAADbmlsZAAKY3JlYXRlZF9hdGQAA25pbGQABWVtYWlsZAADbmlsZAAKZmlyc3RfbmFtZWQAA25pbGQAEWdyb3VwX21lbWJlcnNoaXBzdAAAAARkAA9fX2NhcmRpbmFsaXR5X19kAARtYW55ZAAJX19maWVsZF9fZAARZ3JvdXBfbWVtYmVyc2hpcHNkAAlfX293bmVyX19kACNFbGl4aXIuTWFzc2RyaXZlci5NZWF0c3BhY2UuQWNjb3VudGQACl9fc3RydWN0X19kACFFbGl4aXIuRWN0by5Bc3NvY2lhdGlvbi5Ob3RMb2FkZWRkAAZncm91cHN0AAAABGQAD19fY2FyZGluYWxpdHlfX2QABG1hbnlkAAlfX2ZpZWxkX19kAAZncm91cHNkAAlfX293bmVyX19kACNFbGl4aXIuTWFzc2RyaXZlci5NZWF0c3BhY2UuQWNjb3VudGQACl9fc3RydWN0X19kACFFbGl4aXIuRWN0by5Bc3NvY2lhdGlvbi5Ob3RMb2FkZWRkAAJpZG0AAAAkNmI0ZjRjMjktMDFkMi00YjYyLWIzNzEtMjg3NTM2Y2Y2MDhlZAATaWRlbnRpdHlfc2VydmljZV9pZGQAA25pbGQACWxhc3RfbmFtZWQAA25pbGQADW9yZ2FuaXphdGlvbnN0AAAABGQAD19fY2FyZGluYWxpdHlfX2QABG1hbnlkAAlfX2ZpZWxkX19kAA1vcmdhbml6YXRpb25zZAAJX19vd25lcl9fZAAjRWxpeGlyLk1hc3Nkcml2ZXIuTWVhdHNwYWNlLkFjY291bnRkAApfX3N0cnVjdF9fZAAhRWxpeGlyLkVjdG8uQXNzb2NpYXRpb24uTm90TG9hZGVkZAAKdXBkYXRlZF9hdGQAA25pbG4GAAKbLiSDAWIAAVGA.4_ftk_SunOq4KMF9v0i5k3MRnk3WShk3S7EE2NZQbpU" // nolint:gosec

func NewClient() *graphql.Client {
	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}}
	client := graphql.NewClient("https://api.massdriver.cloud/api/", &c)
	return client
}

func NewSubscriptionClient() *graphql.SubscriptionClient {
	c := http.Client{Transport: &wstransport{underlyingTransport: http.DefaultTransport}}
	client := graphql.NewSubscriptionClient("wss://api.massdriver.cloud/socket/websocket/?vsn=2.0.0")
	client.WithWebSocketOptions(graphql.WebsocketOptions{HTTPClient: &c})
	// TODO not sure if this is necessary for GQL_CONNECTION_INIT
	// client.WithConnectionParams(map[string]interface{}{
	// 	"token": removeMeToken,
	// })
	client.WithLog(func(args ...interface{}) {
		log.Debug().Msgf("%#v", args)
	})
	log.Debug().Msgf("sub client created: %#v", client)
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

// wstransport is a transport for the graphql subscription client because it needs the token set as a url param
// as websockets do not support headers.
type wstransport struct {
	underlyingTransport http.RoundTripper
}

func (t *wstransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO eventually we will use the api key in the token query
	// apiKey := os.Getenv("MASSDRIVER_API_KEY")

	// if apiKey == "" {
	// 	log.Fatal().Msg("MASSDRIVER_API_KEY must be set")
	// }

	// req.URL.Query().Add("token", apiKey)

	// TODO remove this hard coded token
	q := req.URL.Query()
	q.Add("token", removeMeToken)
	req.URL.RawQuery = q.Encode()
	return t.underlyingTransport.RoundTrip(req)
}
