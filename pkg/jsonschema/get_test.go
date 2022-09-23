package jsonschema_test

import (
	"encoding/json"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func TestGet(t *testing.T) {
	type test struct {
		name  string
		input string
		want  jsonschema.Schema
	}
	tests := []test{
		{
			name:  "simple",
			input: "./testdata/schema.json",
			want: jsonschema.Schema{
				Properties: jsonschema.PropertiesMap{
					"firstName": jsonschema.Property{
						Type: "string",
					},
					"lastName": jsonschema.Property{
						Type: "string",
					},
					"age": jsonschema.Property{
						Type:    "integer",
						Minimum: getFloat(0),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := jsonschema.GetJSONSchema(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			gotJSONBytes, _ := json.Marshal(got)
			wantJSONBytes, _ := json.Marshal(tc.want)

			if string(gotJSONBytes) != string(wantJSONBytes) {
				t.Errorf("got %s, want %s", gotJSONBytes, wantJSONBytes)
			}
		})
	}
}

func getFloat(x float64) *float64 {
	return &x
}
