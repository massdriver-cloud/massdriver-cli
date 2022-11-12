package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
)

func TestDeployPreviewEnvironment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Printf("calling this")
		response := map[string]interface{}{}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	os.Setenv("GITHUB_PR", "77")
	defer os.Unsetenv("GITHUB_PR")

	template := `{"hostname": "pr-${GITHUB_PR}.preview.example.com"}`
	previewConfig := strings.NewReader(template)
	environment, err := api.DeployPreviewEnvironment(client, "faux-org-id", "ecomm", previewConfig)

	if err != nil {
		t.Fatal(err)
	}

	got := environment.Config
	want := `{"hostname": "pr-77.preview.example.com"}`

	if got != want {
		t.Errorf("expected: %q, got: %q", got, want)
	}
}
