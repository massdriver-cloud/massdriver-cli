package application_test

import (
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
)

func TestGenerate(t *testing.T) {
	type test struct {
		name     string
		wantPath string
	}
	tests := []test{
		{
			name:     "k8s-app",
			wantPath: "testdata/k8s-app-generate-want",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outDir := t.TempDir()

			templateData := application.TemplateData{
				Name:           "my-test-app",
				TemplateName:   "kubernetes-deployment",
				TemplateSource: "testdata/application-templates",
				OutputDir:      outDir,
			}

			err := application.GenerateFromTemplate(&templateData)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			compareDirs(t, tc.wantPath, templateData.OutputDir)
		})
	}
}
