package application

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"gopkg.in/yaml.v3"
)

const MassdriverApplicationTemplatesRepository = "https://github.com/massdriver-cloud/application-templates"

func Generate(data *TemplateData) error {
	tempDir, err := ioutil.TempDir("/tmp/", "md-app-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	_, cloneErr := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      MassdriverApplicationTemplatesRepository,
		Progress: os.Stdout,
		Depth:    1,
	})

	if cloneErr != nil {
		return ErrCloneFail
	}

	errCopy := copyTemplate(tempDir, data.TemplateName, data.OutputDir)
	if errCopy != nil {
		return ErrCopyFail
	}

	// TODO: only do this for templates w/ a helm chart
	modifyHelmTemplate(tempDir, *data)
	modifyAppYaml(tempDir, *data)

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

func modifyAppYaml(tempDir string, data TemplateData) error {
	appYAML, _ := Parse(tempDir + "/" + data.TemplateName + "/app/app.yaml")
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

	errWrite := ioutil.WriteFile(path.Join(tempDir, data.TemplateName+"/app/", "app.yaml"), appYAMLBytes, common.AllRead|common.UserRW)
	if errWrite != nil {
		return errWrite
	}

	errRename := os.Rename(path.Join(tempDir, data.TemplateName+"/app/app.yaml"), "app/app.yaml")
	if errRename != nil {
		return errRename
	}
	return nil
}

func modifyHelmTemplate(tempDir string, data TemplateData) error {
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
	errWrite := ioutil.WriteFile(path.Join(tempDir, data.TemplateName+"/app/chart", "Chart.yaml"), chartBytes, common.AllRead|common.UserRW)
	if errWrite != nil {
		return errWrite
	}

	errRename := os.Rename(path.Join(tempDir, data.TemplateName+"/app/chart/Chart.yaml"), "app/chart/Chart.yaml")
	if errRename != nil {
		return errRename
	}
	return nil
}
