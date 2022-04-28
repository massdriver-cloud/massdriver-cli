package bundle_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"golang.org/x/mod/sumdb/dirhash"
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
				Ui: map[string]interface{}{
					"ui": "baz",
				},
			},
			apiKey:   "s3cret",
			wantBody: `{"name":"the-bundle","description":"something","type":"bundle","ref":"github.com/some-repo","access":"public","artifacts_schema":"{\"artifacts\":\"foo\"}","connections_schema":"{\"connections\":\"bar\"}","params_schema":"{\"params\":{\"hello\":\"world\"}}","ui_schema":"{\"ui\":\"baz\"}"}`,
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

func TestTarGzipBundle(t *testing.T) {
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
			var got bytes.Buffer
			err := bundle.TarGzipBundle(tc.dirPath, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			// Create a temp dir, write out the archive, then shell out the untar
			testDir := t.TempDir()
			zipOut := path.Join(testDir, "out.tar.gz")
			gotBytes := got.Bytes()
			err = os.WriteFile(zipOut, gotBytes, 0644)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			cmd := exec.Command("tar", "-xzf", zipOut, "-C", testDir)
			err = cmd.Run()
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			wantMD5, err := dirhash.HashDir(path.Dir(tc.dirPath), "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotMD5, err := dirhash.HashDir(path.Join(testDir, "zipdir"), "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(gotMD5) != len(wantMD5) {
				t.Errorf("got %v, want %v", len(gotMD5), len(wantMD5))
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
