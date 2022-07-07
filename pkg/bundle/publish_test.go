package bundle_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func TestPublish(t *testing.T) {
	type test struct {
		name     string
		bundle   bundle.Bundle
		wantBody string
	}
	tests := []test{
		{
			name: "simple",
			bundle: bundle.Bundle{
				Name:        "the-bundle",
				Description: "something",
				Ref:         "github.com/some-repo",
				Type:        "bundle",
				Access:      "public",
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
				UI: map[string]interface{}{
					"ui": "baz",
				},
			},
			wantBody: `{"name":"the-bundle","description":"something","type":"bundle","ref":"github.com/some-repo","access":"public","artifacts_schema":"{\"artifacts\":\"foo\"}","connections_schema":"{\"connections\":\"bar\"}","params_schema":"{\"params\":{\"hello\":\"world\"}}","ui_schema":"{\"ui\":\"baz\"}"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotBody string
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bytes, readErr := ioutil.ReadAll(r.Body)
				if readErr != nil {
					t.Fatalf("%d, unexpected error", readErr)
				}
				gotBody = string(bytes)

				if _, err := w.Write([]byte(`{"upload_location":"https://some.site.test/endpoint"}`)); err != nil {
					t.Fatalf("%d, unexpected error writing upload location to test server", err)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer testServer.Close()

			c := client.NewClient().WithEndpoint(testServer.URL)

			gotResponse, err := tc.bundle.Publish(c)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotBody != tc.wantBody {
				t.Errorf("got %v, want %v", gotBody, tc.wantBody)
			}
			if gotResponse != `https://some.site.test/endpoint` {
				t.Errorf("got %v, want %v", gotResponse, `https://some.site.test/endpoint`)
			}
		})
	}
}

func TestUploadToPresignedS3URL(t *testing.T) {
	type test struct {
		name  string
		bytes []byte
	}
	tests := []test{
		{
			name:  "simple",
			bytes: []byte{1, 2, 3, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotBody []byte
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("%d, unexpected error", err)
				}
				gotBody = bytes

				w.WriteHeader(http.StatusOK)
			}))
			defer testServer.Close()

			err := bundle.UploadToPresignedS3URL(testServer.URL, bytes.NewReader(tc.bytes))
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if string(gotBody) != string(tc.bytes) {
				t.Errorf("got %v, want %v", gotBody, tc.bytes)
			}
		})
	}
}
