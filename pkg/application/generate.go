package application

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"gopkg.in/yaml.v3"
)

func Generate(data *TemplateData) error {
	errCache := cache.GetMassdriverTemplates()
	if errCache != nil {
		return ErrCloneFail
	}
	// TODO: copy template from cache to tmp dir
	tempDir, err := ioutil.TempDir("/tmp/", "md-app-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	errCopy := copyTemplate(cache.TemplateCacheDir, data.TemplateName, data.OutputDir)
	if errCopy != nil {
		return ErrCopyFail
	}

	if errModify := modifyAppYaml(*data); errModify != nil {
		return errModify
	}
	// TODO: only do this for templates w/ a helm chart
	if errModifyHelm := modifyHelmTemplate(*data); errModifyHelm != nil {
		return errModifyHelm
	}

	return nil
}

func copyTemplate(templateDir string, templateName string, outputDir string) error {
	templateFiles, _ := fs.Sub(os.DirFS(templateDir), templateName)

	return fs.WalkDir(templateFiles, ".", func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		outputPath := path.Join(outputDir, filePath)
		if info.IsDir() {
			if filePath == "." {
				return os.MkdirAll(".", common.AllRWX)
			}

			return os.Mkdir(outputPath, common.AllRWX)
		}

		file, readErr := os.ReadFile(templateDir + "/" + templateName + "/" + filePath)
		if readErr != nil {
			return readErr
		}

		writeErr := os.WriteFile(outputPath, file, common.AllRWX)
		if writeErr != nil {
			return writeErr
		}

		return nil
	})
}

func modifyAppYaml(data TemplateData) error {
	appYAML, _ := Parse(data.OutputDir + "/app/app.yaml")
	// TODO: Cory has a PR to change this to title
	appYAML.Name = data.Name
	appYAML.Metadata = Metadata{
		Template: data.TemplateName,
	}
	appYAML.Description = data.Description
	appYAML.Access = data.Access

	appYAMLBytes, err := yaml.Marshal(appYAML)
	if err != nil {
		return err
	}

	errWrite := ioutil.WriteFile(path.Join(data.OutputDir+"/app/", "app.yaml"), appYAMLBytes, common.AllRead|common.UserRW)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func modifyHelmTemplate(data TemplateData) error {
	// regenerate Chart.yaml so k8s labels and selectors are correct
	chart := ChartYAML{
		APIVersion:  "v2",
		Name:        data.Name,
		Description: data.Description,
		Type:        "application",
		Version:     "1.0.0",
	}
	chartBytes, err := yaml.Marshal(chart)
	if err != nil {
		return err
	}
	errWrite := ioutil.WriteFile(path.Join(data.OutputDir, data.TemplateName+"/app/chart", "Chart.yaml"), chartBytes, common.AllRead|common.UserRW)
	if errWrite != nil {
		return errWrite
	}

	return nil
}
