package jsonschema_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func TestWriteDereferencedSchema(t *testing.T) {
	got := bytes.Buffer{}

	err := jsonschema.WriteDereferencedSchema("./testdata/WriteDereferencedSchema/artifact.json", &got, nil)
	if err != nil {
		t.Errorf("Encountered error dereferencing schema: %v", err)
	}

	want, err := os.ReadFile("./testdata/WriteDereferencedSchema/want.json")
	if err != nil {
		t.Errorf("Encountered error dereferencing schema: %v", err)
	}

	if got.String() != string(want) {
		t.Errorf("got %v, want %v", got.String(), string(want))
	}
}
