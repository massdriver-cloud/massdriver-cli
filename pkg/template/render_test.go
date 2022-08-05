package template_test

import (
	"path"
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
		templatesDir string
	}
	tests := []test{
		{
			name:         "my-app",
			description:  "my cool Massdriver app",
			templateName: "kubernetes-deployment",
			templatesDir: "testdata/application-templates",
			wantPath:     "testdata/application-templates-want/kubernetes-deployment",
		},
		{
			name:         "cloud-run-hyperscale-serverless-api",
			description:  "my cool Massdriver app serverless",
			templateName: "cloud-run-api",
			templatesDir: "testdata/application-templates",
			wantPath:     "testdata/application-templates-want/cloud-run-api",
		},
		{
			name:         "my-job",
			description:  "my cool Massdriver cron job",
			templateName: "kubernetes-cronjob",
			templatesDir: "testdata/application-templates",
			wantPath:     "testdata/application-templates-want/kubernetes-cronjob",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			templateData := template.Data{
				Name:           tc.name,
				Description:    tc.description,
				TemplateName:   tc.templateName,
				TemplateSource: tc.templatesDir,
				OutputDir:      t.TempDir(),
				CloudProvider:  "gcp",
				Dependencies: map[string]string{
					"massdriver/draft-node_0": "massdriver/draft-node",
				},
			}
			templateDir := path.Join(tc.templatesDir, tc.templateName)
			err := template.RenderDirectory(templateDir, &templateData)
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
