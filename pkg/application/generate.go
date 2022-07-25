package application

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"gopkg.in/yaml.v3"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Generate(data *TemplateData) error {
	templates, _ := cache.ApplicationTemplates()
	if !contains(templates, data.TemplateName) {
		return fmt.Errorf("template '%s' not found, try `mass app templates refresh`", data.TemplateName)
	}

	errCopy := copyTemplate(cache.AppTemplateCacheDir(), data.TemplateName, data.OutputDir)
	if errCopy != nil {
		return ErrCopyFail
	}

	// TODO: only do this for templates w/ a helm chart
	if errModifyHelm := modifyHelmTemplate(*data); errModifyHelm != nil {
		return errModifyHelm
	}
	if errModify := modifyAppYaml(*data); errModify != nil {
		return errModify
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
	appYAML, _ := Parse(data.OutputDir + "/" + data.TemplateName + "/app/app.yaml")
	// TODO: Cory has a PR to change this to title
	appYAML.Title = data.Name
	appYAML.Metadata = Metadata{
		Template: data.TemplateName,
	}
	appYAML.Description = data.Description
	appYAML.Access = data.Access

	appYAMLBytes, err := yaml.Marshal(appYAML)
	if err != nil {
		return err
	}

	errWrite := ioutil.WriteFile(path.Join(data.OutputDir, data.TemplateName+"/app/", "app.yaml"), appYAMLBytes, common.AllRead|common.UserRW)
	if errWrite != nil {
		return errWrite
	}

	errRename := os.Rename(path.Join(data.OutputDir, data.TemplateName+"/app/app.yaml"), "app/app.yaml")
	if errRename != nil {
		return errRename
	}
	return nil
}

func modifyHelmTemplate(data TemplateData) error {
	// regenerate Chart.yaml to match their config
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

	errRename := os.Rename(path.Join(data.OutputDir, data.TemplateName+"/app/chart/Chart.yaml"), "app/chart/Chart.yaml")
	if errRename != nil {
		return errRename
	}
	return nil
}
