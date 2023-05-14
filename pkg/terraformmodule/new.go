package terraformmodule

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
//go:embed templates/* templates/terraform-module/.pre-commit-config.yaml templates/terraform-module/.gitignore templates/terraform-module/_providers.tf templates/terraform-module/_variables.tf templates/massdriver-application/.pre-commit-config.yaml templates/massdriver-application/.gitignore templates/massdriver-application/_providers.tf templates/massdriver-application/_variables.tf
var templatesFs embed.FS

func Generate(data *template.Data) error {
	// We are replicating the liquid template logic here, from our repo-manager
	// {{- $cloud_prefix := $bundle_name | regexMatch "^[a-z]+-" | strings.Trim "-" -}}
	// {{- $repo_name := printf "massdriver-cloud/%s" $bundle_name -}}
	// {{- $repo_encoded := $repo_name | regexp.Replace "/" "%2F" -}}

	// add cloud prefix
	r := regexp.MustCompile("^[a-z]+-")
	data.CloudPrefix = strings.Trim(r.FindString(data.Name), "-")

	// add repo name
	data.RepoName = fmt.Sprintf("massdriver-cloud/%s", data.Name)
	data.RepoNameEncoded = strings.ReplaceAll(data.RepoName, "/", "%2F")

	templateDir, _ := fs.Sub(fs.FS(templatesFs), path.Join("templates", data.TemplateName))
	return template.RenderEmbededDirectory(templateDir, data)
}

func GenerateApplication(data *template.Data) error {
	// We are replicating the liquid template logic here, from our repo-manager
	// {{- $cloud_prefix := $bundle_name | regexMatch "^[a-z]+-" | strings.Trim "-" -}}
	// {{- $repo_name := printf "massdriver-cloud/%s" $bundle_name -}}
	// {{- $repo_encoded := $repo_name | regexp.Replace "/" "%2F" -}}

	// add cloud prefix
	r := regexp.MustCompile("^[a-z]+-")
	data.CloudPrefix = strings.Trim(r.FindString(data.Name), "-")

	// add repo name
	data.RepoName = fmt.Sprintf("massdriver-cloud/%s", data.Name)
	data.RepoNameEncoded = strings.ReplaceAll(data.RepoName, "/", "%2F")

	templateDir, _ := fs.Sub(fs.FS(templatesFs), path.Join("templates", data.TemplateName))
	return template.RenderEmbededDirectory(templateDir, data)
}
