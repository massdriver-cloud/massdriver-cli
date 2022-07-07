package terraform_test

import (
	"encoding/json"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/terraform"
)

type test struct {
	name  string
	input jsonschema.Property
	want  terraform.TFVariable
}

func TestNewTFVariable(t *testing.T) {
	tests := []test{
		{
			name:  "scalars",
			input: jsonschema.Property{Type: "number"},
			want:  terraform.TFRequiredVariable{Type: "number"},
		},
		{
			name:  "list of scalars",
			input: jsonschema.Property{Type: "array", Items: &jsonschema.Property{Type: "string"}},
			want:  terraform.TFRequiredVariable{Type: "any"},
		},
		{
			name:  "list of any",
			input: jsonschema.Property{Type: "array", Items: &jsonschema.Property{}},
			want:  terraform.TFRequiredVariable{Type: "any"},
		},
		{
			name:  "maps",
			input: jsonschema.Property{Type: "object", AdditionalProperties: true},
			want:  terraform.TFRequiredVariable{Type: "any"},
		},
		{
			name:  "object w/ scalars",
			input: jsonschema.Property{Type: "object", Properties: jsonschema.PropertiesMap{"street_number": jsonschema.Property{Type: "number"}, "street_name": jsonschema.Property{Type: "string"}}},
			want:  terraform.TFRequiredVariable{Type: "any"},
		},
		{
			name: "complex objects",
			input: jsonschema.Property{
				Type: "object",
				Properties: jsonschema.PropertiesMap{
					"name": jsonschema.Property{Type: "string"},
					"children": jsonschema.Property{
						Type: "array",
						Items: &jsonschema.Property{
							Type: "object",
							Properties: jsonschema.PropertiesMap{
								"name": jsonschema.Property{Type: "string"},
							},
						},
					},
				},
			},
			want: terraform.TFRequiredVariable{Type: "any"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := terraform.NewTFVariable(tc.input, true)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestTFVariableFileJSONEncoding(t *testing.T) {
	type test struct {
		name  string
		input terraform.TFVariableFile
		want  string
	}

	tests := []test{
		{
			name:  "A single variable",
			input: terraform.TFVariableFile{Variable: map[string]terraform.TFVariable{"name": terraform.TFRequiredVariable{Type: "string"}}},
			want:  `{"variable":{"name":{"type":"string"}}}`,
		},
		{
			name:  "Multiple variables",
			input: terraform.TFVariableFile{Variable: map[string]terraform.TFVariable{"name": terraform.TFRequiredVariable{Type: "string"}, "age": terraform.TFRequiredVariable{Type: "number"}}},
			want:  `{"variable":{"age":{"type":"number"},"name":{"type":"string"}}}`,
		},
		{
			name:  "An optional variable",
			input: terraform.TFVariableFile{Variable: map[string]terraform.TFVariable{"name": terraform.TFOptionalVariable{Type: "string"}}},
			want:  `{"variable":{"name":{"type":"string","default":null}}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bytes, _ := json.Marshal(tc.input)
			got := string(bytes)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
