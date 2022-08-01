package cmd

import (
	"errors"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
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

	// TODO: mo-betta, apps / bunbdles / etc should share metadata
	app, err := application.Parse("massdriver.yaml")
	if err != nil {
		return err
	}

	switch app.Type {
	case "application":
		return application.Build(c, output)
	case "bundle":
		return bundle.Build(c, output)
	default:
		return errors.New("unknown type")
	}
}
