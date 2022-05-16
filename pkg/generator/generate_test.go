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
		Name:        "aws-vpc",
		Access:      "Private",
		Description: "a vpc",
		Type:        "bundle",
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

	err := generator.Generate(&bundleData)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	templatePath := fmt.Sprintf("%s/%s", bundleData.BundleDir, bundleData.Name)

	bundleYamlPath := fmt.Sprintf("%s/massdriver.yaml", templatePath)
	expectedContent := "name: aws-vpc"

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
