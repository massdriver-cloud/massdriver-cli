package cmd

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	masscmd "github.com/massdriver-cloud/massdriver-cli/pkg/cmd"
	"github.com/massdriver-cloud/massdriver-cli/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:     "preview",
	Aliases: []string{"pv"},
	Short:   "Preview Environments",
	Long:    ``,
}

var previewParamsPath string
var previewCiContextPath string
var previewInitCmd = &cobra.Command{
	Use:   "init project_slug",
	Short: "Generates a preview environment config file for your project.",
	RunE:  runPreviewInit,
	Args:  cobra.ExactArgs(1),
}

var previewDeployCmd = &cobra.Command{
	Use:   "deploy project_slug",
	Short: "Deploys a preview environment in your project.",
	RunE:  runPreviewDeploy,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(previewCmd)

	previewCmd.AddCommand(previewInitCmd)
	previewInitCmd.Flags().StringVarP(&previewParamsPath, "output", "o", "./preview.json", "Output path for preview environment params file. This file supports bash interpolation and can be manually edited or programatically modified during CI.")

	previewCmd.AddCommand(previewDeployCmd)
	previewDeployCmd.Flags().StringVarP(&previewParamsPath, "params", "p", "./preview.json", "Path to preview params file. This file supports bash interpolation.")
	previewDeployCmd.Flags().StringVarP(&previewCiContextPath, "ci-context", "c", "", "Path to GitHub Actions event.json")
}

func runPreviewDeploy(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)
	c := config.Get()

	projectSlugOrID := args[0]

	client := api2.NewClient(c.APIKey)
	environment, err := masscmd.DeployPreviewEnvironment(client, c.OrgID, projectSlugOrID, previewParamsPath, previewCiContextPath)

	if err != nil {
		log.Error().Err(err).Msg("Failed to deploy environment")
		return err
	}

	log.Info().Str("id", environment.ID).Str("Config", previewParamsPath).Msgf("Preview environment deploying.")
	return nil
}

func runPreviewInit(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)
	c := config.Get()
	projectSlugOrID := args[0]

	err := masscmd.InitializePreview(c, projectSlugOrID, previewParamsPath)

	if err != nil {
		return err
	}

	log.Info().Str("id", projectSlugOrID).Msgf("Preview environment default parameters output to %s", previewParamsPath)
	return nil
}
