package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog/log"
)

func (app *Application) Build(c *client.MassdriverClient, output string, appPath string) error {
	log.Info().Msg("Building application...")

	// TODO: be better
	bun, err := bundle.Parse(appPath, nil)
	if err != nil {
		return err
	}

	if errBuild := bun.Build(c, output); errBuild != nil {
		return errBuild
	}
	log.Info().Msg("Application built")

	return nil
}
