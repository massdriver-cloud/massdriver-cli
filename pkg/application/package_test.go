package application_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"golang.org/x/mod/sumdb/dirhash"
)

func TestPackage(t *testing.T) {
	type test struct {
		name            string
		applicationPath string
		wantPath        string
	}
	tests := []test{
		{
			name:            "simple",
			applicationPath: "testdata/appsimple/massdriver.yaml",
			wantPath:        "testdata/simple",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got bytes.Buffer

			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				schemaFile := r.URL.Path
				switch schemaFile {
				case "/artifact-definitions/massdriver/kubernetes-cluster":
					if _, err := w.Write([]byte(`{"kube":"cluster"}`)); err != nil {
						t.Errorf("Encountered error writing kube cluster: %v", err)
					}
				case "/artifact-definitions/massdriver/aws-iam-role":
					if _, err := w.Write([]byte(`{"aws":"authentication"}`)); err != nil {
						t.Errorf("Encountered error writing aws iam role: %v", err)
					}
				case "/artifact-definitions/massdriver/gcp-service-account":
					if _, err := w.Write([]byte(`{"gcp":"authentication"}`)); err != nil {
						t.Errorf("Encountered error writing gcp service account: %v", err)
					}
				case "/artifact-definitions/massdriver/azure-service-principal":
					if _, err := w.Write([]byte(`{"azure":"authentication"}`)); err != nil {
						t.Errorf("Encountered error writing azure service principal: %v", err)
					}
				default:
					t.Fatalf("unknown schema: %v", schemaFile)
				}
			}))
			defer testServer.Close()

			c := client.NewClient().WithBaseURL(testServer.URL)

			// Create a temp dir, write out the archive, then shell out to untar
			testDir := t.TempDir()

			_, err := application.Package(tc.applicationPath, c, testDir, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			wantMD5, err := dirhash.HashDir(tc.wantPath, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotMD5, err := dirhash.HashDir(testDir, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotMD5 != wantMD5 {
				t.Errorf("got %v, want %v", gotMD5, wantMD5)
			}
		})
	}
}
