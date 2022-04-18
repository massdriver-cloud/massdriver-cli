package jsonschema_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func TestWriteDereferencedSchema(t *testing.T) {
	dir, err := ioutil.TempDir("", "xo-artifacts")
	if err != nil {
		t.Errorf("%d, unexpected error", err)
	}

	err = jsonschema.WriteDereferencedSchema("./testdata/WriteDereferencedSchema/aws-authentication.json", dir)
	if err != nil {
		t.Errorf("Encountered error dereferencing schema: %v", err)
	}

	gotDir, _ := os.ReadDir(dir)

	got := []string{}
	for _, dirEntry := range gotDir {
		got = append(got, dirEntry.Name())
	}

	want := []string{"aws-authentication.dereferenced.json"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	defer os.RemoveAll(dir)
}
