package jsonschema_test

import (
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func TestValidateJSONDocument(t *testing.T) {
	got, _ := jsonschema.Validate("testdata/schema.json", "testdata/valid-document.json")
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestInvalidateJSONDocument(t *testing.T) {
	got, _ := jsonschema.Validate("testdata/valid-schema.json", "testdata/invalid-document.json")
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}
