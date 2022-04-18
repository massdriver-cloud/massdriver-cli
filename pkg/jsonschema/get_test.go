package jsonschema_test

import (
	"fmt"
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
						Type: "integer",
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

			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
