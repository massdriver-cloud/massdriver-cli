package api2_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

func TestListCredentials(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"artifacts": map[string]interface{}{
					"items": []map[string]interface{}{
						{
							"id":   "uuid1",
							"name": "artifact1",
						},
						{
							"id":   "uuid2",
							"name": "artifact2",
						},
					},
				},
			},
		}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	credentials, err := api2.ListCredentials(client, "faux-org-id", "massdriver/aws-iam-role")

	if err != nil {
		t.Fatal(err)
	}

	got := len(credentials)
	want := 2

	if got != want {
		t.Errorf("got %d credentials, wanted %d", got, want)
	}
}
