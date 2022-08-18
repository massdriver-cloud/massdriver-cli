package bundle_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"gopkg.in/yaml.v3"
)

func TestTransformations(t *testing.T) {
	type testData struct {
		name           string
		schemaPath     string
		transformation func(map[string]interface{}) error
		expected       map[string]interface{}
	}
	tests := []testData{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tc.schemaPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := map[string]interface{}{}

			err = yaml.Unmarshal(data, &got)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			err = bundle.ApplyTransformations(got, []func(map[string]interface{}) error{tc.transformation})
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if fmt.Sprint(got) != fmt.Sprint(tc.expected) {
				t.Errorf("got %v, want %v", got, tc.expected)
			}
		})
	}
}
