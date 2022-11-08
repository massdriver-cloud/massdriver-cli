package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj"},
	Short:   "Project management",
	Long:    ``,
}

var projectGetDefaultParams = &cobra.Command{
	Use:   "get-default-params",
	Short: "Gets default parameters for all manifests in project.",
	RunE:  runProjectGetDefaultParams,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectGetDefaultParams)
}

func runProjectGetDefaultParams(cmd *cobra.Command, args []string) error {
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

	return writePreviewEnvFile(project, "preview.json")
}

func writePreviewEnvFile(project *api.Project, outfile string) error {
	log.Info().Str("id", project.ID).Str("slug", project.Slug).Msgf("Preview environment default parameters output to %s", outfile)

	previewConf, err := json.MarshalIndent(project.DefaultParams, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outfile, previewConf, 0600)
	return err
}
