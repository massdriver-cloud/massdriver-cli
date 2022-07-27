package application

import (
	"fmt"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func Build(appPath string, workingDir string, c *client.MassdriverClient) error {
	app, parseErr := Parse(appPath)
	if parseErr != nil {
		return parseErr
	}
	_, err := app.ConvertToBundle()
	if err != nil {
		return fmt.Errorf("could not convert app to bundle: %w", err)
	}
	err = bundle.Build(workingDir, c)
	if err != nil {
		return err
	}
	return nil
}
