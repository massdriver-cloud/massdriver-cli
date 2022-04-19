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
	var bundle, _ = bundle.ParseBundle("./testdata/bundle.Build/bundle.yaml")
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
	var b, _ = bundle.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	var inputIo bytes.Buffer

	bundle.GenerateSchema(b.Params, b.Metadata("params"), &inputIo)
	var gotJson = &map[string]interface{}{}
	_ = json.Unmarshal(inputIo.Bytes(), gotJson)

	wantBytes, _ := ioutil.ReadFile("./testdata/bundle.Build/schema-params.json")
	var wantJson = &map[string]interface{}{}
	_ = json.Unmarshal(wantBytes, wantJson)

	if !reflect.DeepEqual(gotJson, wantJson) {
		t.Errorf("got %v, want %v", gotJson, wantJson)
	}
}
