package bundle

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
	"github.com/rs/zerolog/log"
)

const TerraformTemplateName = "terraform-module"

func Generate(data *template.Data) error {
	data.TemplateName = TerraformTemplateName
	log.Info().Msgf("Generating bundle from template '%s'", data.TemplateName)
	templates, err := cache.ApplicationTemplates()
	if err != nil {
		return err
	}

	if !common.Contains(templates, data.TemplateName) {
		return fmt.Errorf("template '%s' not found, try `mass app templates refresh`", data.TemplateName)
	}

	// add cloud prefix
	r := regexp.MustCompile("^[a-z]+-")
	data.CloudPrefix = strings.Trim(r.FindString(data.Name), "-")

	// add repo name
	data.RepoName = fmt.Sprintf("massdriver-cloud/%s", data.Name)
	data.RepoNameEncoded = strings.ReplaceAll(data.RepoName, "/", "%2F")

	data.TemplateSource = cache.AppTemplateCacheDir()
	templatePath := path.Join(data.TemplateSource, data.TemplateName)
	return template.RenderDirectory(templatePath, data)
}
