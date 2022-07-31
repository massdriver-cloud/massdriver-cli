package application

import (
	"fmt"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
	"github.com/rs/zerolog/log"
)

func GenerateFromTemplate(data *template.Data) error {
	log.Info().Msgf("Generating application from template '%s'", data.TemplateName)
	templates, err := cache.ApplicationTemplates()
	if err != nil {
		return err
	}

	if !common.Contains(templates, data.TemplateName) {
		return fmt.Errorf("template '%s' not found, try `mass app templates refresh`", data.TemplateName)
	}

	data.TemplateSource = cache.AppTemplateCacheDir()
	templatePath := path.Join(data.TemplateSource, data.TemplateName)
	return template.RenderDirectory(templatePath, data)
}
