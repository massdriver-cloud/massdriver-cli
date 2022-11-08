package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
)

func TestGetProject(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"project": map[string]interface{}{
					"id":            "00000000-1111-2222-3333-444444444444",
					"defaultParams": `{"network":{"region": "us-west-2"}, "cluster":{"region": "us-west-2"}}`,
				},
			},
		}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	project, err := api.GetProject(client, "faux-org-id", "ecomm")

	if err != nil {
		t.Fatal(err)
	}

	got := project.ID
	want := "00000000-1111-2222-3333-444444444444"

	if got != want {
		t.Errorf("got project.ID: %q, want: %q", got, want)
	}
}
