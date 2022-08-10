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
			},
			want: &bundle.Bundle{
				Name:        "app-name",
				Description: "description",
				Ref:         "github.com/some-repo",
				Access:      "public",
				Type:        "application",
				Params: map[string]interface{}{
					"params": map[string]interface{}{
						"hello": "world",
					},
				},
				Connections: map[string]interface{}{
					"required": []string{"another-field"},
					"properties": map[string]interface{}{
						"another-field": map[string]interface{}{"$ref": "bar"},
						"some-field":    map[string]interface{}{"$ref": "foo"},
					},
				},
				Artifacts: map[string]interface{}{
					"properties": map[string]interface{}{},
				},
				UI: map[string]interface{}{
					"ui:order": []interface{}{"*"},
				},
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
