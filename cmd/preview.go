package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:     "preview",
	Aliases: []string{"pv"},
	Short:   "Preview Environments",
	Long:    ``,
}

var previewConfigPath string
var previewInitCmd = &cobra.Command{
	Use:   "init project_slug",
	Short: "Generate an environment params template for creating preview environments.",
	RunE:  runPreviewInit,
	Args:  cobra.ExactArgs(1),
}

var previewCreateCmd = &cobra.Command{
	Use:   "create project_slug",
	Short: "Creates and deploys a new preview environment.",
	RunE:  runPreviewCreate,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(previewCmd)

	previewCmd.AddCommand(previewInitCmd)
	previewInitCmd.Flags().StringVarP(&previewConfigPath, "output", "o", "./preview.json", "Output path for preview environment parameters file")

	previewCmd.AddCommand(previewCreateCmd)
	previewCreateCmd.Flags().StringVarP(&previewConfigPath, "config", "c", "./preview.json", "Path to preview environment parameters file")
}

func runPreviewCreate(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	projectSlugOrId := args[0]

	orgID := os.Getenv("MASSDRIVER_ORG_ID")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_ORG_ID must be set")
	}

	previewConfig, err := os.Open(previewConfigPath)

	if err != nil {
		return err
	}

	client := api.NewClient()
	environment, err := api.CreatePreviewEnvironment(client, orgID, projectSlugOrId, previewConfig)

	_ = environment

	if err != nil {
		log.Error().Err(err).Msg("Failed to get project")
		return err
	}

	return nil
}

func runPreviewInit(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	projectSlugOrId := args[0]

	orgID := os.Getenv("MASSDRIVER_ORG_ID")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_ORG_ID must be set")
	}

	client := api.NewClient()
	project, err := api.GetProject(client, orgID, projectSlugOrId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get project")
		return err
	}

	return writePreviewConfigFile(project, previewConfigPath)
}

func writePreviewConfigFile(project *api.Project, path string) error {
	log.Info().Str("id", project.ID).Str("slug", project.Slug).Msgf("Preview environment default parameters output to %s", path)

	previewConf, err := json.MarshalIndent(project.DefaultParams, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, previewConf, 0600)
	return err
}
