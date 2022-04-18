package application_test

import (
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
)

func TestReportDeploymentStatus(t *testing.T) {

	type testData struct {
		name       string
		configPath string
		want       bool
	}
	tests := []testData{
		{
			name:       "Test Valid Config",
			configPath: "testdata/valid.json",
			want:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := application.Lint(tc.configPath, "file://testdata/schema-application.json")
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if got != tc.want {
				t.Fatalf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}
