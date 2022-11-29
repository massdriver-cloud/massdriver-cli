package bundle_test

import (
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestParseBundle(t *testing.T) {
	var got, _ = bundle.Parse("./testdata/massdriver.yaml", nil)
	var want = bundle.Bundle{
		Schema:      "draft-07",
		Type:        "infrastructure",
		Name:        "aws-vpc",
		Description: "Something",
		Access:      "public",
		Artifacts:   map[string]interface{}{},
		Params: map[string]interface{}{
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":  "string",
					"title": "Name",
				},
				"age": map[string]interface{}{
					"type":  "integer",
					"title": "Age",
				},
			},
			"required": []interface{}{
				"name",
			},
		},
		Connections: map[string]interface{}{
			"required": []interface{}{
				"default",
			},
			"properties": map[string]interface{}{
				"default": map[string]interface{}{
					"type":  "string",
					"title": "Default credential",
				},
			},
		},
		UI: map[string]interface{}{},
	}

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("got %v, want %v", *got, want)
	}
}

func TestParseBundleWithAccessOverride(t *testing.T) {
	overrides := map[string]interface{}{
		"access": "private",
	}
	var bundle, _ = bundle.Parse("./testdata/massdriver.yaml", overrides)
	got := bundle.Access
	want := "private"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
