package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog/log"
)

func Build(b *bundle.Bundle, c *client.MassdriverClient, output string) error {
	log.Info().Msg("Building application...")
	// app-specific build logic goes here //
	if errBuild := b.Build(c, output); errBuild != nil {
		return errBuild
	}
	log.Info().Msg("Application built")

	return nil
}
