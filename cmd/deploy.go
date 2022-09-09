/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a configured package",
	Long:  ``,
	Args:  cobra.ExactArgs(1),

	RunE: runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)
	name := args[0]

	orgID := os.Getenv("MASSDRIVER_ORG_ID")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_ORG_ID must be set")
	}

	client := api.NewClient()
	subClient := api.NewSubscriptionClient()
	deployment, err := api.DeployPackage(client, subClient, orgID, name)

	if err != nil {
		if deployment != nil {
			log.Fatal().Err(err).Str("deploymentId", deployment.ID).Msg("Deployment failed")
		} else {
			log.Fatal().Err(err).Msg("Deployment failed")
		}
		return err
	}

	log.Info().Str("deploymentId", deployment.ID).Msgf("Deployment %s", strings.ToLower(deployment.Status))
	return nil
}
