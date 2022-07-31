package bundle

import (
	"embed"
	"io/fs"

	// "text/template"

	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

// note to all: option in go 1.18 will load hidden files so we dont have to include `cp` instructions in readme for pre-commit.
//go:embed templates/* templates/terraform/.pre-commit-config.yaml templates/terraform/.gitignore templates/terraform/.github/workflows/build.yaml templates/terraform/.github/workflows/publish.yaml templates/terraform/src/_artifacts.tf templates/terraform/src/_providers.tf
var templatesFs embed.FS

// type TemplateData struct {
// 	Name        string
// 	Description string
// 	Access      string
// 	Type        string
// 	OutputDir   string
// }

func Generate(data *template.Data) error {
	templateDir, _ := fs.Sub(fs.FS(templatesFs), "templates/terraform")
	return template.RenderEmbededDirectory(templateDir, data)
}
