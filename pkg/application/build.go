package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

var configFile = "massdriver.yaml"

func (app *Application) Build(c *client.MassdriverClient, output string) error {
	// TODO: be better
	bun, err := bundle.Parse(configFile, nil)
	if err != nil {
		return err
	}

	if errBuild := bun.Build(c, output); errBuild != nil {
		return errBuild
	}

	return nil
}
