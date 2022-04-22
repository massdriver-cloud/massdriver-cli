package bundle_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

				w.Write([]byte(`{"upload_location":"https://some.site.test/endpoint"}`))
				w.WriteHeader(http.StatusOK)
			}))
			defer testServer.Close()

			bundle.MASSDRIVER_URL = testServer.URL

			gotResponse, err := tc.bundle.Publish(tc.apiKey)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotBody != tc.wantBody {
				t.Errorf("got %v, want %v", gotBody, tc.wantBody)
			}
			if gotResponse != `https://some.site.test/endpoint` {
				t.Errorf("got %v, want %v", gotResponse, `https://some.site.test/endpoint`)
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

func TestTarGzipDirectory(t *testing.T) {
	type test struct {
		name     string
		dirPath  string
		wantFile string
	}
	tests := []test{
		{
			name:     "simple",
			dirPath:  "testdata/zipdir/bundle.yaml",
			wantFile: "testdata/zipdir.tar.gz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantBytes, err := ioutil.ReadFile(tc.wantFile)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			var got bytes.Buffer

			err = bundle.TarGzipDirectory(tc.dirPath, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotBytes := got.Bytes()
			err = os.WriteFile("/tmp/dat1.tar.gz", gotBytes, 0644)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(gotBytes) != len(wantBytes) {
				t.Errorf("got %v, want %v", len(gotBytes), len(wantBytes))
			}
		})
	}
}
