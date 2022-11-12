package api2_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

func TestGetDeployment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"deployment": map[string]interface{}{
					"id":     "uuid1",
					"status": "PROVISIONING",
				},
			},
		}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	deployment, err := api2.GetDeployment(client, "faux-org-id", "uuid1")

	if err != nil {
		t.Fatal(err)
	}

	got := deployment.Status
	want := "PROVISIONING"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
