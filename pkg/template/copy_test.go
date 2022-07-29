package template_test

import (
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
	"golang.org/x/mod/sumdb/dirhash"
)

func TestAppFromTemplate(t *testing.T) {
	type test struct {
		name         string
		description  string
		templateName string
		wantPath     string
		templateDir  string
	}
	tests := []test{
		{
			name:         "my-app",
			description:  "my cool Massdriver app",
			templateName: "kubernetes-deployment",
			templateDir:  "testdata/application-templates",
			wantPath:     "testdata/application-templates-want/kubernetes-deployment",
		},
		{
			name:         "cloud-run-hyperscale-serverless-api",
			description:  "my cool Massdriver app serverless",
			templateName: "cloud-run-api",
			templateDir:  "testdata/application-templates",
			wantPath:     "testdata/application-templates-want/cloud-run-api",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			templateData := template.Data{
				Name:         tc.name,
				Description:  tc.description,
				TemplateName: tc.templateName,
				// OutputDir:    t.TempDir(),

				OutputDir: "_local-test",
			}
			err := template.Copy(tc.templateDir, &templateData)
			if err != nil {
				t.Errorf("unexpected error copying template: %v", err)
			}

			wantMD5, err := dirhash.HashDir(tc.wantPath, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			gotMD5, err := dirhash.HashDir(templateData.OutputDir, "", dirhash.DefaultHash)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if gotMD5 != wantMD5 {
				// TODO: need to print out which file is different
				t.Errorf("got %v, want %v", gotMD5, wantMD5)
			}
		})
	}
}
