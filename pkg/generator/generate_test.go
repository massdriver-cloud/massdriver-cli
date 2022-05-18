package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/generator"
)

func TestGenerate(t *testing.T) {
	bundleData := generator.TemplateData{
		Name:        "aws-vpc",
		Access:      "Private",
		Description: "a vpc",
		Type:        "bundle",
		OutputDir:   "./testdata/bundle",
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

	os.Mkdir(bundleData.OutputDir, 0777)

	err := generator.Generate(&bundleData)
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

	t.Cleanup(func() {
		os.RemoveAll(bundleData.OutputDir)
	})
}
