package api2_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

func TestDeployPreviewEnvironment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// TODO: This isn't necessary just here for my sanity trying to debug JSON! w/ genqlient
		var params map[string]interface{}
		json.NewDecoder(req.Body).Decode(&params)
		input := params["variables"].(map[string]interface{})["input"]
		ciContextStr := input.(map[string]interface{})["ciContext"]
		ciContext := map[string]interface{}{}
		json.Unmarshal([]byte(fmt.Sprintf("%v", ciContextStr)), &ciContext)
		pullRequest := ciContext["pull_request"]
		prNum := int((pullRequest.(map[string]interface{})["number"]).(float64))
		slug := fmt.Sprintf("p%d", prNum)

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"deployPreviewEnvironment": map[string]interface{}{
					"result": map[string]interface{}{
						"id":   "env-uuid1",
						"slug": slug,
					},
					"successful": true,
				},
			},
		}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	confMap := map[string]interface{}{
		"network": map[string]interface{}{
			"cidr": "10.0.0.0/16",
		},
		"cluster": map[string]interface{}{
			"maxNodes": 10,
		},
	}
	ciContext := map[string]interface{}{
		"pull_request": map[string]interface{}{
			"title":  "First commit!",
			"number": 69,
		},
	}

	credentials := []api2.Credential{}

	environment, err := api2.DeployPreviewEnvironment(client, "faux-org-id", "faux-project-id", credentials, confMap, ciContext)

	// TODO test interpolation GITHUB_PR, etc
	if err != nil {
		t.Fatal(err)
	}

	got := environment.ID
	want := "env-uuid1"

	if got != want {
		t.Errorf("got %s , wanted %s", got, want)
	}
}
