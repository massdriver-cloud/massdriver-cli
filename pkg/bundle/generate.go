package bundle

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

// note to all: option in go 1.18 will load hidden files so we dont have to include `cp` instructions in readme for pre-commit.
//go:embed templates/* templates/terraform/.pre-commit-config.yaml templates/terraform/.gitignore templates/terraform/.github/workflows/build.yaml templates/terraform/.github/workflows/publish.yaml templates/terraform/src/_artifacts.tf templates/terraform/src/_providers.tf
var templatesFs embed.FS

type TemplateData struct {
	Name        string
	Description string
	Access      string
	Type        string
	OutputDir   string
}

func Generate(data *TemplateData) error {
	templateFiles, _ := fs.Sub(fs.FS(templatesFs), "templates/terraform")

	err := fs.WalkDir(templateFiles, ".", func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		outputPath := path.Join(data.OutputDir, filePath)
		if info.IsDir() {
			if filePath == "." {
				return os.MkdirAll(data.OutputDir, common.AllRWX)
			}

			return os.Mkdir(outputPath, common.AllRWX)
		}

		var tmpl *template.Template
		var outputFile *os.File
		tmpl, _ = template.ParseFS(templateFiles, filePath)
		outputFile, err = os.Create(outputPath)

		if err != nil {
			return err
		}

		defer outputFile.Close()
		return tmpl.Execute(outputFile, data)
	})

	return err
}
