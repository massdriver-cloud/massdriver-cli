package application

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func Build(c *client.MassdriverClient, output string) error {
	app, err := Parse("massdriver.yaml")
	if err != nil {
		return err
	}
	_, errBun := app.ConvertToBundle()
	if errBun != nil {
		return errBun
	}

	if errBuild := bundle.Build(c, output); errBuild != nil {
		return errBuild
	}

	return nil
}
