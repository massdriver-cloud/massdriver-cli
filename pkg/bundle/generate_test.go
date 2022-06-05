package bundle_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
)

func TestGenerate(t *testing.T) {
	bundleData := bundle.TemplateData{
		Name:        "aws-vpc",
		Access:      "Private",
		Description: "a vpc",
		Type:        "bundle",
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

	terraformPath := fmt.Sprintf("%s/src", templatePath)
	mainTFPath := fmt.Sprintf("%s/main.tf", terraformPath)
	expectedContent = regexp.MustCompile("random_pet")
	assertFileCreatedAndContainsText(t, mainTFPath, expectedContent)
}
