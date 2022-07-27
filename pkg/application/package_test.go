package application_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func TestPackage(t *testing.T) {
	type test struct {
		name            string
		applicationPath string
		wantPath        string
	}
	tests := []test{
		{
			name:            "k8s-app",
			applicationPath: "testdata/k8s-app-generate-want/app.yaml",
			wantPath:        "testdata/k8s-app-package-want",
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

			c := client.NewClient().WithEndpoint(testServer.URL)

			// Create a temp dir, write out the archive, then shell out to untar
			testDir := t.TempDir()

			_, err := application.PackageApplication(tc.applicationPath, c, testDir, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			compareDirs(t, tc.wantPath, testDir)
		})
	}
}
