package application

import (
	"fmt"

	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
	"github.com/rs/zerolog/log"
)

func GenerateFromTemplate(data *template.Data) error {
	log.Info().Msgf("Generating application from template %v", data)
	templates, _ := cache.ApplicationTemplates()
	if !common.Contains(templates, data.TemplateName) {
		return fmt.Errorf("template '%s' not found, try `mass app templates refresh`", data.TemplateName)
	}
	source := data.TemplateSource
	if source == "" {
		source = cache.AppTemplateCacheDir()
	}

	// TODO: use template.Data higher up the call chain
	tmplData := &template.Data{
		TemplateName: data.TemplateName,
		Name:         data.Name,
		Description:  data.Description,
	}

	errCopy := template.Copy(source, tmplData)
	if errCopy != nil {
		log.Err(errCopy).Msg("error copying template")
		return errCopy
	}

	return nil
}
