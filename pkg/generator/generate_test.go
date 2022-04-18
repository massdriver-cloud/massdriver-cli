package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/massdriver-cloud/massdriver-cli/pkg/generator"
)

func TestGenerate(t *testing.T) {
	//TODO: We should be mocking the filesystem here.
	//The testing/testFS package isn't quite there yet and afero although cool seems like it has implications
	//for the broader application.
	bundleData := generator.TemplateData{
		Name:        "AWS VPC",
		Type:        "aws-vpc",
		Access:      "Private",
		Description: "a vpc",
		TemplateDir: "./testdata/templates",
		BundleDir:   "./testdata/bundle",
	}

	assertFileCreatedAndContainsText := func(t testing.TB, filename, expectedContent string) {
		t.Helper()
		content, err := ioutil.ReadFile(filename)

		if err != nil {
			t.Errorf("Failed to create %s", filename)
		}

		if !strings.Contains(string(content), expectedContent) {
			t.Errorf("Data failed to render in template %s", filename)
		}
	}

	os.Mkdir(bundleData.BundleDir, 0777)

	generator.Generate(bundleData)

	templatePath := fmt.Sprintf("%s/%s", bundleData.BundleDir, bundleData.Type)

	bundleYamlPath := fmt.Sprintf("%s/bundle.yaml", templatePath)
	expectedContent := "title: AWS VPC"

	assertFileCreatedAndContainsText(t, bundleYamlPath, expectedContent)

	readmePath := fmt.Sprintf("%s/README.md", templatePath)
	expectedContent = "a vpc"

	assertFileCreatedAndContainsText(t, readmePath, expectedContent)

	terraformPath := fmt.Sprintf("%s/terraform", templatePath)
	mainTFPath := fmt.Sprintf("%s/main.tf", terraformPath)
	expectedContent = "random_pet"

	assertFileCreatedAndContainsText(t, mainTFPath, expectedContent)

	t.Cleanup(func() {
		os.RemoveAll(bundleData.BundleDir)
	})
}
