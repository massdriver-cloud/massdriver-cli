package application_test

import (
	"errors"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestLintEnvs(t *testing.T) {
	type test struct {
		name string
		app  *bundle.Bundle
		want map[string]string
		err  error
	}
	tests := []test{
		{
			name: "params, connections, secrets working",
			app: &bundle.Bundle{
				App: &bundle.AppBlock{
					Envs: map[string]string{
						"FOO":        ".params.foo",
						"INTEGER":    ".params.int",
						"CONNECTION": ".connections.connection1",
						"SECRET":     ".secrets.shh",
					},
					Secrets: map[string]bundle.Secret{
						"shh": {},
					},
				},
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"foo": map[string]interface{}{
							"type":  "string",
							"const": "bar",
						},
						"int": map[string]interface{}{
							"type":  "integer",
							"const": 4,
						},
					},
				},
				Connections: map[string]interface{}{
					"properties": map[string]interface{}{
						"connection1": map[string]interface{}{
							"type":  "string",
							"const": "whatever",
						},
					},
				},
			},
			want: map[string]string{
				"FOO":        "bar",
				"INTEGER":    "4",
				"CONNECTION": "whatever",
				"SECRET":     "some-secret-value",
			},
			err: nil,
		},
		{
			name: "error on missing data",
			app: &bundle.Bundle{
				App: &bundle.AppBlock{
					Envs: map[string]string{
						"FOO": ".params.foo",
					},
					Secrets: map[string]bundle.Secret{},
				},
				Params:      map[string]interface{}{},
				Connections: map[string]interface{}{},
			},
			want: map[string]string{},
			err:  errors.New("The jq query for environment variable FOO didn't produce a result"),
		},
		{
			name: "error on invalid jq syntax",
			app: &bundle.Bundle{
				App: &bundle.AppBlock{
					Envs: map[string]string{
						"FOO": "laksdjf",
					},
					Secrets: map[string]bundle.Secret{},
				},
				Params:      map[string]interface{}{},
				Connections: map[string]interface{}{},
			},
			want: map[string]string{},
			err:  errors.New("The jq query for environment variable FOO produced an error: function not defined: laksdjf/0"),
		},
		{
			name: "error on multiple values",
			app: &bundle.Bundle{
				App: &bundle.AppBlock{
					Envs: map[string]string{
						"FOO": ".params.array[]",
					},
					Secrets: map[string]bundle.Secret{},
				},
				Params: map[string]interface{}{
					"properties": map[string]interface{}{
						"array": map[string]interface{}{
							"type":     "array",
							"minItems": 2,
							"items": map[string]interface{}{
								"type": "integer",
							},
						},
					},
				},
				Connections: map[string]interface{}{},
			},
			want: map[string]string{},
			err:  errors.New("The jq query for environment variable FOO produced multiple values, which isn't supported"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := application.LintEnvs(tc.app)
			if tc.err != nil {
				if err == nil {
					t.Errorf("expected an error, got nil")
				} else if tc.err.Error() != err.Error() {
					t.Errorf("got %v, want %v", err.Error(), tc.err.Error())
				}
			} else if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(got) != len(tc.want) {
				t.Errorf("got %v, want %v", len(got), len(tc.want))
			}
			for key, wantValue := range tc.want {
				gotValue, ok := got[key]
				if !ok {
					t.Errorf("got %v, want %v", got, tc.want)
				}
				if gotValue != wantValue {
					t.Errorf("got %v, want %v", gotValue, wantValue)
				}
			}
		})
	}
}
