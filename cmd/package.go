package cmd

import (
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "Pacakge development tools",
	Long:    ``,
}

var packageConfigureCmd = &cobra.Command{
	Use:     `configure <project>-<target>-<manifest>`,
	Aliases: []string{"config"},
	Short:   "Configure a package in a Massdriver target",
	Long:    ``,
	RunE:    runPackageConfigure,
	Args:    cobra.ExactArgs(1),
}

var PackageConfigureJSONPath string

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.AddCommand(packageConfigureCmd)
	packageConfigureCmd.Flags().StringVar(&PackageConfigureJSONPath, "f", "", "Path to JSON or YAML file to use for package configuration")
}

func runPackageConfigure(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	orgID := os.Getenv("MASSDRIVER_ORG_ID")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_ORG_ID must be set")
	}

	client := api.NewClient()
	_, err := api.ConfigurePackage(client, orgID, args[0], PackageConfigureJSONPath)
	return err
}
