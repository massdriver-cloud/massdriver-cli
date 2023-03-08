package template_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/mod/sumdb/dirhash"
)

func outputDir(t *testing.T) string {
	return t.TempDir()
}

func TestAppFromTemplate(t *testing.T) {
	type test struct {
		name         string
		description  string
		templateName string
		wantPath     string
		templatesDir string
		connections  map[string]interface{}
	}
	tests := []test{
		{
			name:         "renders-successfully",
			description:  "Renders an application config & subdirectories",
			templateName: "renders-successfully",
			templatesDir: "testdata/application-templates/",
			wantPath:     "testdata/application-templates-want/renders-successfully",
		},
		{
			name:         "renders-connections",
			description:  "Renders selectected dependencies as Connections",
			templateName: "renders-connections",
			templatesDir: "testdata/application-templates/",
			wantPath:     "testdata/application-templates-want/renders-connections",
			connections: map[string]interface{}{
				"draft_node": "massdriver/draft-node",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			templateData := template.Data{
				Name:           tc.name,
				Description:    tc.description,
				TemplateName:   tc.templateName,
				TemplateSource: tc.templatesDir,
				OutputDir:      outputDir(t),
				// OutputDir:   "_local-test",
				Connections: tc.connections,
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
				t.Errorf("got %v, want %v", gotMD5, wantMD5)
				walkAndCompare(tc.wantPath, templateData.OutputDir)
			}
		})
	}
}

func walkAndCompare(wantDir string, gotDir string) {
	err := filepath.Walk(wantDir,
		func(path string, info os.FileInfo, err error) error {
			isDir, _ := isDirectory(path)

			if isDir {
				return nil
			}

			relativeFilePath := strings.TrimPrefix(path, wantDir)
			gotFilePath := filepath.Join(gotDir, relativeFilePath)

			if err != nil {
				return err
			}

			gotText, _ := readFile(gotFilePath)
			wantText, _ := readFile(path)
			if gotText != wantText {
				fmt.Printf("File did not render correctly: %s\n", path)
				fmt.Printf("Comparing (want) %s and (got) %s\n", path, gotFilePath)
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(wantText, gotText, false)
				fmt.Println(dmp.DiffToDelta(diffs))
				fmt.Printf("==== Got ====\n%s\n==== End Got ====\n", gotText)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
