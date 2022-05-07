package application_test

import (
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestConvertToBundle(t *testing.T) {
	type test struct {
		name string
		app  *application.Application
		want *bundle.Bundle
	}
	tests := []test{
		{
			name: "simple",
			app: &application.Application{
				Name:        "app-name",
				Description: "description",
				Ref:         "github.com/some-repo",
				Access:      "public",
				Dependencies: []application.ApplicationDependencies{
					{
						Type:     "foo",
						Field:    "some-field",
						Required: false,
						Env:      []application.ApplicationDependenciesEnvs{},
					},
					{
						Type:     "bar",
						Field:    "another-field",
						Required: true,
						Env:      []application.ApplicationDependenciesEnvs{},
					},
				},
				Params: map[string]interface{}{
					"params": map[string]string{
						"hello": "world",
					},
				},
			},
			want: &bundle.Bundle{
				Name:        "app-name",
				Description: "description",
				Ref:         "github.com/some-repo",
				Access:      "public",
				Type:        "application",
				Params: map[string]interface{}{
					"params": map[string]string{
						"hello": "world",
					},
				},
				Connections: map[string]interface{}{
					"required": []string{"kubernetes-cluster", "cloud-authentication", "another-field"},
					"properties": map[string]interface{}{
						"another-field":        map[string]string{"$ref": "bar"},
						"cloud-authentication": map[string]string{"$ref": "massdriver/cloud-authentication"},
						"some-field":           map[string]string{"$ref": "foo"},
						"kubernetes-cluster":   map[string]string{"$ref": "massdriver/kubernetes-cluster"},
					},
				},
				Artifacts: map[string]interface{}{
					"required": []string{"kubernetes-application"},
					"properties": map[string]interface{}{
						"kubernetes-application": map[string]string{"$ref": "massdriver/kubernetes-application"},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 	schemaFile := r.URL.Path
			// 	switch schemaFile {
			// 	case "/artifact-definitions/massdriver/kubernetes-cluster":
			// 		w.Write([]byte(`{"kube":"cluster"}`))
			// 	case "/artifact-definitions/massdriver/cloud-authentication":
			// 		w.Write([]byte(`{"cloud":"authentication"}`))
			// 	case "/artifact-definitions/foo":
			// 		w.Write([]byte(`{"hello":"world"}`))
			// 	case "/artifact-definitions/bar":
			// 		w.Write([]byte(`{"lol":"rofl"}`))
			// 	default:
			// 		t.Fatalf("unknown schema: %v", schemaFile)
			// 	}
			// }))
			// defer testServer.Close()

			// c := client.NewClient().WithEndpoint(testServer.URL)

			got := tc.app.ConvertToBundle()

			if !reflect.DeepEqual(*got, *tc.want) {
				t.Errorf("got %v, want %v", *got, *tc.want)
			}
		})
	}
}
