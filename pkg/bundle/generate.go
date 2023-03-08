package bundle

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

// note to all: option in go 1.18 will load hidden files so we dont have to include `cp` instructions in readme for pre-commit.
//
//go:embed templates/* templates/terraform/.pre-commit-config.yaml templates/terraform/.gitignore templates/terraform/.github/workflows/build.yaml templates/terraform/.github/workflows/publish.yaml templates/terraform/src/_artifacts.tf templates/terraform/src/_providers.tf
var templatesFs embed.FS

func Generate(data *template.Data) error {
	// add cloud prefix
	r := regexp.MustCompile("^[a-z]+-")
	data.CloudPrefix = strings.Trim(r.FindString(data.Name), "-")

	// add repo name
	data.RepoName = fmt.Sprintf("massdriver-cloud/%s", data.Name)
	data.RepoNameEncoded = strings.ReplaceAll(data.RepoName, "/", "%2F")

	templateDir, _ := fs.Sub(fs.FS(templatesFs), path.Join("templates", data.TemplateName))
	return template.RenderEmbededDirectory(templateDir, data)
}
