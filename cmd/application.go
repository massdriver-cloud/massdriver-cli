/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:     "application",
	Aliases: []string{"app"},
	Short:   "Application development tools",
	Long:    ``,
}

var applicationGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Deprecated: Generates a new application template",
	RunE:    runApplicationGenerate,
}

var applicationNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"new"},
	Short:   "Creates a new application from a template",
	RunE:    runApplicationNew,
}

var applicationPublishCmd = &cobra.Command{
	Use:          "publish [Path to app.yaml]",
	Short:        "Publish an application to Massdriver",
	Args:         cobra.ExactArgs(1),
	RunE:         runApplicationPublish,
	SilenceUsage: true,
}

var applicationTemplatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"plates"},
	Short:   "Lists available application templates",
	RunE:    runApplicationTemplates,
}

var templatesRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshes local copy of application templates",
	RunE:  runTemplatesRefresh,
}

func init() {
	rootCmd.AddCommand(applicationCmd)

	applicationCmd.AddCommand(applicationGenerateCmd)
	applicationCmd.AddCommand(applicationNewCmd)
	applicationCmd.AddCommand(applicationPublishCmd)
	applicationCmd.AddCommand(applicationTemplatesCmd)
	applicationTemplatesCmd.AddCommand(templatesRefreshCmd)
}

func runApplicationGenerate(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templateData := application.TemplateData{}

	err := application.RunPrompt(&templateData)
	if err != nil {
		return err
	}

	err = application.Generate(&templateData)
	if err != nil {
		return err
	}

	return nil
}

func runApplicationNew(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templateData := application.TemplateData{}

	err := application.RunPrompt(&templateData)
	if err != nil {
		return err
	}

	err = application.GenerateFromTemplate(&templateData)
	if err != nil {
		return err
	}

	return nil
}

func runApplicationPublish(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	var err error
	appPath := args[0]

	c := client.NewClient()

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithAPIKey(apiKey)
	}

	// Create a temporary working directory
	workingDir, err := os.MkdirTemp("", "application")
	if err != nil {
		return (err)
	}
	defer os.RemoveAll(workingDir)

	var buf bytes.Buffer
	b, err := application.PackageApplication(appPath, c, workingDir, &buf)
	if err != nil {
		return err
	}

	uploadURL, err := b.Publish(c)
	if err != nil {
		return err
	}

	err = bundle.UploadToPresignedS3URL(uploadURL, &buf)
	if err != nil {
		return err
	}

	fmt.Println("Application published successfully!")

	return nil
}

func runApplicationTemplates(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templates, err := cache.ApplicationTemplates()
	if err != nil {
		return err
	}
	log.Info().Msgf("Application templates:\n  %s", strings.Join(templates, "\n  "))

	return nil
}

func runTemplatesRefresh(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	if err := cache.RefreshAppTemplates(); err != nil {
		return err
	}
	log.Info().Msg("Application templates refreshed successfully.")

	return nil
}
