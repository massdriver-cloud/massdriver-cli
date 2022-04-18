package terraform

import (
	"encoding/json"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"reflect"
	"testing"
)

type test struct {
	name  string
	input jsonschema.Property
	want  TFVariable
}

func TestNewTFVariable(t *testing.T) {
	tests := []test{
		{
			name:  "scalars",
			input: jsonschema.Property{Type: "number"},
			want:  TFRequiredVariable{Type: "number"},
		},
		{
			name:  "list of scalars",
			input: jsonschema.Property{Type: "array", Items: &jsonschema.Property{Type: "string"}},
			want:  TFRequiredVariable{Type: "any"},
		},
		{
			name:  "list of any",
			input: jsonschema.Property{Type: "array", Items: &jsonschema.Property{}},
			want:  TFRequiredVariable{Type: "any"},
		},
		{
			name:  "maps",
			input: jsonschema.Property{Type: "object", AdditionalProperties: true},
			want:  TFRequiredVariable{Type: "any"},
		},
		{
			name:  "object w/ scalars",
			input: jsonschema.Property{Type: "object", Properties: jsonschema.PropertiesMap{"street_number": jsonschema.Property{Type: "number"}, "street_name": jsonschema.Property{Type: "string"}}},
			want:  TFRequiredVariable{Type: "any"},
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
			want: TFRequiredVariable{Type: "any"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NewTFVariable(tc.input, true)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestTFVariableFileJSONEncoding(t *testing.T) {
	type test struct {
		name  string
		input TFVariableFile
		want  string
	}

	tests := []test{
		{
			name:  "A single variable",
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFRequiredVariable{Type: "string"}}},
			want:  `{"variable":{"name":{"type":"string"}}}`,
		},
		{
			name:  "Multiple variables",
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFRequiredVariable{Type: "string"}, "age": TFRequiredVariable{Type: "number"}}},
			want:  `{"variable":{"age":{"type":"number"},"name":{"type":"string"}}}`,
		},
		{
			name:  "An optional variable",
			input: TFVariableFile{Variable: map[string]TFVariable{"name": TFOptionalVariable{Type: "string"}}},
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
