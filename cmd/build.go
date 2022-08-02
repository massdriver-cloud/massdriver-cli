package cmd

import (
	"errors"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:  "build",
	Long: ``,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	// TODO: app/bundle build directories
	output := "."

	// TODO: apps / bundles / etc should share metadata
	app, err := application.Parse("massdriver.yaml")
	if err != nil {
		return err
	}

	switch app.Type {
	case "application":
		return app.Build(c, output)
	case "bundle":
		bun, errParse := bundle.Parse("massdriver.yaml", nil)
		if errParse != nil {
			return errParse
		}
		return bun.Build(c, output)
	default:
		return errors.New("unknown type")
	}
}
