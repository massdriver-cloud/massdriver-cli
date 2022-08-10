package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func (app *Application) Build(c *client.MassdriverClient, output string) error {
	bundle, err := app.ConvertToBundle()
	if err != nil {
		return err
	}

	if errBuild := bundle.Build(c, output); errBuild != nil {
		return errBuild
	}

	return nil
}
