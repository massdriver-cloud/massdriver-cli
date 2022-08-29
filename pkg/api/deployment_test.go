//
package api_test

import (
	"net/http"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
)

func TestGetDeployment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mustWrite(w, `{"data":{"deployment":{"id":"foo", "status": "RUNNING"}}}`)
	})
	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	deployment, err := api.GetDeployment(client, "faux-org-id", "foo")

	if err != nil {
		t.Fatal(err)
	}

	if got, want := deployment.Status, "RUNNING"; got != want {
		t.Errorf("got deployment.ID: %q, want: %q", got, want)
	}
}
