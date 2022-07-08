package bundle_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestGenerateSchemas(t *testing.T) {
	var bundle, _ = bundle.Parse("./testdata/bundle.Build/massdriver.yaml", nil)
	_ = bundle.GenerateSchemas("./tmp/")

	gotDir, _ := os.ReadDir("./tmp")
	got := []string{}

	for _, dirEntry := range gotDir {
		got = append(got, dirEntry.Name())
	}

	want := []string{"schema-artifacts.json", "schema-connections.json", "schema-params.json", "schema-ui.json"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	defer os.RemoveAll("./tmp")
}

func TestGenerateSchema(t *testing.T) {
	var b, _ = bundle.Parse("./testdata/bundle.Build/massdriver.yaml", nil)
	var inputIo bytes.Buffer

	if err := bundle.GenerateSchema(b.Params, b.Metadata("params"), &inputIo); err != nil {
		t.Errorf("Encountered error generating schema: %v", err)
	}
	var gotJSON = &map[string]interface{}{}
	_ = json.Unmarshal(inputIo.Bytes(), gotJSON)

	wantBytes, _ := ioutil.ReadFile("./testdata/bundle.Build/schema-params.json")
	var wantJSON = &map[string]interface{}{}
	_ = json.Unmarshal(wantBytes, wantJSON)

	if !reflect.DeepEqual(gotJSON, wantJSON) {
		t.Errorf("got %v, want %v", gotJSON, wantJSON)
	}
}
