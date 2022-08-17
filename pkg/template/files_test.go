package template_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	gotmpl "text/template"

	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

func TestWriteTemplateToFile(t *testing.T) {
	type test struct {
		name     string
		template string
		fileName string
		wantFile string
		data     template.Data
	}
	tests := []test{
		{
			name: "Test name",
			template: `# Massdriver Application Template
# Template: {{mass .TemplateName mass}}
title: {{mass .Name mass}}
description: {{mass .Description mass}}
`,
			data: template.Data{
				Name:         "App Name",
				Description:  "App Description",
				TemplateName: "app-template",
				OutputDir:    "app-diretory",
			},
			fileName: "here-file",
			wantFile: "testdata/templates/yaml.yaml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, _ := gotmpl.New("tmpl").Delims(template.OpenPattern, template.ClosePattern).Parse(tc.template)
			filePath := path.Join(t.TempDir(), tc.fileName)
			err := template.WriteToFile(filePath, tmpl, &tc.data)
			if err != nil {
				t.Errorf("unexpected error copying template: %v", err)
			}

			gotFile, _ := ioutil.ReadFile(filePath)
			wantFile, _ := ioutil.ReadFile(tc.wantFile)

			if !bytes.Equal(gotFile, wantFile) {
				t.Errorf("\ngot \n%s\n\nwant \n%s", gotFile, wantFile)
			}
		})
	}
}
