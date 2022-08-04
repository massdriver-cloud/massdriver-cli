package bundle_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

func TestGenerate(t *testing.T) {
	bundleData := template.Data{
		Name:         "aws-vpc",
		Access:       "Private",
		Description:  "a vpc",
		Type:         "bundle",
		TemplateName: "terraform",
	}

	assertFileCreatedAndContainsText := func(t testing.TB, filename string, expectedPattern *regexp.Regexp) {
		t.Helper()
		content, err := ioutil.ReadFile(filename)

		if err != nil {
			t.Errorf("Failed to create %s", filename)
		}

		if !expectedPattern.MatchString(string(content)) {
			t.Errorf("Data failed to render in template %s", filename)
		}
	}

	testDir := t.TempDir()
	bundleData.OutputDir = testDir

	err := bundle.Generate(&bundleData)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	templatePath := bundleData.OutputDir

	bundleYamlPath := fmt.Sprintf("%s/massdriver.yaml", templatePath)
	expectedContent := regexp.MustCompile("name.*aws-vpc")

	assertFileCreatedAndContainsText(t, bundleYamlPath, expectedContent)

	readmePath := fmt.Sprintf("%s/README.md", templatePath)
	expectedContent = regexp.MustCompile("# aws-vpc")
	assertFileCreatedAndContainsText(t, readmePath, expectedContent)

	srcPath := fmt.Sprintf("%s/src", templatePath)
	mainTFPath := fmt.Sprintf("%s/main.tf", srcPath)
	expectedContent = regexp.MustCompile("random_pet")
	assertFileCreatedAndContainsText(t, mainTFPath, expectedContent)

	validationsJSONPath := fmt.Sprintf("%s/validations.json", srcPath)
	expectedContent = regexp.MustCompile("do_not_delete")
	assertFileCreatedAndContainsText(t, validationsJSONPath, expectedContent)
}
