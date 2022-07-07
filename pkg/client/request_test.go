package client_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func TestToHTTPRequest(t *testing.T) {
	type test struct {
		name    string
		request client.Request
		want    http.Request
	}
	tests := []test{
		{
			name: "simple",
			request: client.Request{
				Method: "GET",
				Path:   "some/path",
				Body:   strings.NewReader("some data"),
			},
			want: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "https",
					Host:   "api.massdriver.cloud",
					Path:   "/some/path",
				},
				Body: io.NopCloser(strings.NewReader("some data")),
				Header: http.Header{
					"X-Md-Api-Key": []string{"apikey"},
					"Content-Type": []string{"application/json"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			c := client.NewClient().WithAPIKey("apikey")
			got, err := tc.request.toHTTPRequest(ctx, c)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if got.Method != tc.want.Method {
				t.Errorf("got %v, want %v", got.Method, tc.want.Method)
			}
			if got.URL.String() != tc.want.URL.String() {
				t.Errorf("got %v, want %v", got.URL.String(), tc.want.URL.String())
			}
			var gotBody []byte
			var wantBody []byte
			if _, err := got.Body.Read(gotBody); err != nil {
				t.Errorf("could not read got body %v", string(gotBody))
			}
			tc.want.Body.Read(wantBody)
			if string(gotBody) != string(wantBody) {
				t.Errorf("got %v, want %v", string(gotBody), string(wantBody))
			}

			if len(got.Header) != len(tc.want.Header) {
				t.Errorf("got %v, want %v", len(got.Header), len(tc.want.Header))
			}
			for k, v := range tc.want.Header {
				if len(got.Header[k]) != len(v) {
					t.Errorf("got %v, want %v", len(got.Header[k]), len(v))
				}
				for i := range v {
					if len(got.Header[k][i]) != len(v[i]) {
						t.Errorf("got %v, want %v", got.Header[k][i], len(v[i]))
					}
				}
			}
		})
	}
}
