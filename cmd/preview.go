package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
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

	projectSlugOrId := args[0]

	previewConfig, err := os.Open(previewParamsPath)

	if err != nil {
		return err
	}

	client := api.NewClient()
	environment, err := api.DeployPreviewEnvironment(client, c.OrgId, projectSlugOrId, previewConfig)

	_ = environment

	if err != nil {
		log.Error().Err(err).Msg("Failed to get project")
		return err
	}

	return nil
}

func runPreviewInit(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)
	c := config.Get()
	projectSlugOrId := args[0]

	client := api.NewClient()
	project, err := api.GetProject(client, c.OrgId, projectSlugOrId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get project")
		return err
	}

	client2 := api2.NewClient(c.APIKey)
	selectedArtifacts, err := masscmd.InitializePreview(client2, c.OrgId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get artifacts")
		return err
	}

	fmt.Printf("Artifacts: %v\n", selectedArtifacts)
	fmt.Printf("Default Params: %v\n", project.DefaultParams)

	conf := map[string]interface{}{
		"artifacts":     selectedArtifacts,
		"packageParams": project.DefaultParams,
	}

	log.Info().Str("id", project.ID).Str("slug", project.Slug).Msgf("Preview environment default parameters output to %s", previewParamsPath)
	return writePreviewConfigFile(conf, previewParamsPath)
}

func writePreviewConfigFile(conf map[string]interface{}, path string) error {

	previewConf, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, previewConf, 0600)
	return err
}
