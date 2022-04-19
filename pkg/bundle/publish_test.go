package bundle_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestPublish(t *testing.T) {
	type test struct {
		name        string
		bundle      bundle.Bundle
		apiKey      string
		wantBody    string
		wantHeaders map[string][]string
	}
	tests := []test{
		{
			name: "simple",
			bundle: bundle.Bundle{
				Uuid:   "deadbeef-0000",
				Title:  "The Bundle",
				Type:   "bundle-type",
				Access: "public",
				Artifacts: map[string]interface{}{
					"artifacts": "foo",
				},
				Connections: map[string]interface{}{
					"connections": "bar",
				},
				Params: map[string]interface{}{
					"params": map[string]string{
						"hello": "world",
					},
				},
				Ui: map[string]interface{}{
					"ui": "baz",
				},
			},
			apiKey:   "s3cret",
			wantBody: `{"name":"The Bundle","ref":"bundle-type","id":"deadbeef-0000","access":"public","artifacts_schema":{"artifacts":"foo"},"connections_schema":{"connections":"bar"},"params_schema":{"params":{"hello":"world"}},"ui_schema":{"ui":"baz"}}`,
			wantHeaders: map[string][]string{
				"X-Md-Api-Key": {"s3cret"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotBody string
			var gotHeaders map[string][]string
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("%d, unexpected error", err)
				}
				gotBody = string(bytes)
				gotHeaders = r.Header
				w.WriteHeader(http.StatusOK)
			}))
			defer testServer.Close()

			bundle.MASSDRIVER_URL = testServer.URL

			err := tc.bundle.Publish(tc.apiKey)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotBody != tc.wantBody {
				t.Errorf("got %v, want %v", gotBody, tc.wantBody)
			}
			for key, wantValue := range tc.wantHeaders {
				if gotValue, ok := gotHeaders[key]; ok {
					if len(gotValue) != len(wantValue) {
						t.Errorf("got %v, want %v", gotValue, wantValue)
					}
					for i, v := range wantValue {
						if v != gotValue[i] {
							t.Errorf("got %v, want %v", gotValue, wantValue)
						}
					}
				} else {
					t.Errorf("missing header: %v", key)
				}
			}
		})
	}
}
