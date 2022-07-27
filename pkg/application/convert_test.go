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
				Dependencies: map[string]application.Dependencies{
					"some-field": {
						Type:     "foo",
						Required: false,
						Envs:     []application.DependenciesEnvs{},
					},
					"another-field": {
						Type:     "bar",
						Required: true,
						Envs:     []application.DependenciesEnvs{},
					},
				},
				Params: map[string]interface{}{
					"params": map[string]interface{}{
						"hello": "world",
					},
				},
				UI: map[string]interface{}{
					"ui:order": []interface{}{"*"},
				},
			},
			want: &bundle.Bundle{
				Name:        "app-name",
				Description: "description",
				Ref:         "github.com/some-repo",
				Access:      "public",
				Type:        "application",
				Steps: []bundle.Step{
					{
						Path:        "src",
						Provisioner: "terraform",
					},
				},
				Params: map[string]interface{}{
					"params": map[string]interface{}{
						"hello": "world",
					},
				},
				Connections: map[string]interface{}{
					"required": []string{"kubernetes_cluster", "another-field"},
					"properties": map[string]interface{}{
						"another-field": map[string]interface{}{"$ref": "bar"},
						// "cloud-authentication": map[string]interface{}{
						// 	"oneOf": []interface{}{
						// 		map[string]interface{}{"$ref": "massdriver/aws-iam-role"},
						// 		map[string]interface{}{"$ref": "massdriver/azure-service-principal"},
						// 		map[string]interface{}{"$ref": "massdriver/gcp-service-account"},
						// 	},
						// },
						"some-field":           map[string]interface{}{"$ref": "foo"},
						"kubernetes_cluster":   map[string]interface{}{"$ref": "massdriver/kubernetes-cluster"},
						"aws_authentication":   map[string]interface{}{"$ref": "massdriver/aws-iam-role"},
						"azure_authentication": map[string]interface{}{"$ref": "massdriver/azure-service-principal"},
						"gcp_authentication":   map[string]interface{}{"$ref": "massdriver/gcp-service-account"},
					},
				},
				Artifacts: map[string]interface{}{
					"properties": map[string]interface{}{},
				},
				UI: map[string]interface{}{
					"ui:order": []interface{}{"*"},
				},
				// 	"required": []string{"kubernetes-application"},
				// 	"properties": map[string]interface{}{
				// 		"kubernetes-application": map[string]interface{}{"$ref": "massdriver/kubernetes-application"},
				// 	},
				// },
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.app.ConvertToBundle()
			if err != nil {
				t.Errorf("unexpected error converting app to bundle: %v", err)
			}

			if !reflect.DeepEqual(*got, *tc.want) {
				t.Errorf("got %v, want %v", *got, *tc.want)
			}
		})
	}
}
