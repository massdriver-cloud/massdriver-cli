package api2_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

func TestDockerRegistryToken(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc(APIURL, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"data": map[string]map[string]string{
				"containerRepository": {
					"token":   "bogustoken",
					"repoUri": "massdriveruswest.pkg.docker.dev",
				},
			},
		}

		data, _ := json.Marshal(response)
		mustWrite(w, string(data))
	})

	client := graphql.NewClient(APIURL, &http.Client{Transport: localRoundTripper{handler: mux}})

	got, err := api2.GetContainerRepository(client, "artifactId", "orgId", "westus", "massdriver/test-image")

	if err != nil {
		t.Fatal(err)
	}

	want := &api2.ContainerRepository{
		Token:         "bogustoken",
		RepositoryUri: "massdriveruswest.pkg.docker.dev",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Wanted %v but got %v", want, got)
	}
}
