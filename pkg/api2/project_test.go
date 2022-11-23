package api2_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

func TestGetProject(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"project": map[string]interface{}{
					"id":   "uuid1",
					"slug": "sluggy",
					"defaultParams": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	project, err := api2.GetProject(client, "faux-org-id", "sluggy")

	if err != nil {
		t.Fatal(err)
	}

	got := project.Slug
	want := "sluggy"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
