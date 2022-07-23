/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"os"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:              "application",
	Aliases:          []string{"app"},
	Short:            "Application development tools",
	Long:             ``,
	TraverseChildren: true,
}

var applicationGenerateCmd = &cobra.Command{
	Use:              "generate",
	Aliases:          []string{"gen"},
	Short:            "Generates a new application template",
	RunE:             runApplicationGenerate,
	TraverseChildren: true,
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

	applicationAddFlags(applicationGenerateCmd)
	applicationCmd.AddCommand(applicationGenerateCmd)
	applicationCmd.AddCommand(applicationPublishCmd)
	applicationCmd.AddCommand(applicationTemplatesCmd)
	applicationTemplatesCmd.AddCommand(templatesRefreshCmd)
}

var nameDefault = ""
var descriptionDefault = "placeholder description, written to app.yaml"
var templateDefault = ""
var accessDefault = "private"

func applicationAddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("name", "n", nameDefault, "What's it called?")
	cmd.PersistentFlags().StringP("template", "t", templateDefault, "Which application-template to use?")
	// access is not exposed in the CLI, but can be manually set in the generated app.yaml
	// cmd.PersistentFlags().StringP("access", "a", accessDefault, "public or priviate")
}

func runApplicationGenerate(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	name, nameErr := cmd.Flags().GetString("name")
	if nameErr != nil {
		log.Err(nameErr).Msg("Failed to generate an application")
		return nil
	}
	template, templateErr := cmd.Flags().GetString("template")
	if templateErr != nil {
		log.Err(templateErr).Msg("Failed to generate an application")
		return nil
	}

	templateData := application.TemplateData{
		Name:         name,
		Description:  descriptionDefault,
		TemplateName: template,
		Access:       accessDefault,
	}
	errPrompt := application.RunPrompt(&templateData)
	if errPrompt != nil {
		log.Err(errPrompt).Msg("Failed to generate an application")
		return nil
	}

	errGenerate := application.Generate(&templateData)
	if errGenerate != nil {
		log.Err(errGenerate).Msg("Failed to generate an application")
		return nil
	}

	return nil
}

func runApplicationPublish(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c, err := initClient(cmd)
	if err != nil {
		return err
	}
	appPath := args[0]

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

	log.Info().Msg("Application published successfully!")

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
