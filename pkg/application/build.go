package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog/log"
)

func (app *Application) Build(c *client.MassdriverClient, output string) error {
	log.Info().Msg("Building application...")

	// we were trying to avoid any conversion from app to bundle etc
	// this is cheap until the build function is consolidated
	b := app.AsBundle()
	if errBuild := b.Build(c, output); errBuild != nil {
		return errBuild
	}
	log.Info().Msg("Application built")

	return nil
}
