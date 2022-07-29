package template

import (
	"io/fs"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"gopkg.in/osteele/liquid.v1"
)

type Data struct {
	Name           string
	Description    string
	Access         string
	Chart          string
	Location       string
	TemplateName   string
	TemplateSource string
	OutputDir      string
}

func CopyTemplate(templateDir string, data *Data) error {
	templateName := data.TemplateName
	outputDir := data.OutputDir
	templateFiles, _ := fs.Sub(os.DirFS(templateDir), templateName)

	engine := liquid.NewEngine().Delims("<md", "md>", "<tag", "tag>")
	bindings := map[string]interface{}{
		"data": map[string]string{
			"title":       data.Name,
			"description": data.Description,
			"template":    data.TemplateName,
			"cli-version": "v0.1.0",
		},
	}

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
		template := string(file)
		out, err := engine.ParseAndRenderString(template, bindings)
		if err != nil {
			return err
		}

		writeErr := os.WriteFile(outputPath, []byte(out), common.AllRWX)
		if writeErr != nil {
			return writeErr
		}

		return nil
	})
}
