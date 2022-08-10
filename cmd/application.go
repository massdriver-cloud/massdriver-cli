/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
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

var applicationBuildCmd = &cobra.Command{
	Use:  "build",
	RunE: runApplicationBuild,
}

var applicationNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new application from a template",
	RunE:  runApplicationNew,
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
	Aliases: []string{"tmpl"},
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

	applicationCmd.AddCommand(applicationBuildCmd)
	applicationCmd.AddCommand(applicationNewCmd)
	applicationCmd.AddCommand(applicationPublishCmd)
	applicationCmd.AddCommand(applicationTemplatesCmd)
	applicationTemplatesCmd.AddCommand(templatesRefreshCmd)
}

func checkIsApplication(app *application.Application) error {
	if app.Type != "application" {
		return fmt.Errorf("mass app build can only be used with application type massdriver.yaml")
	}
	return nil
}

func runApplicationBuild(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}
	// TODO: app/bundle build directories
	output := "."

	app, err := application.Parse("massdriver.yaml")
	if err != nil {
		return err
	}
	if errType := checkIsApplication(app); errType != nil {
		return errType
	}
	return app.Build(c, output)
}

func runApplicationNew(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templateData := template.Data{
		Access: "private",
		// TODO: unify bundle build and app build outputDir logic and support
		OutputDir: ".",
	}

	err := application.RunPromptNew(&templateData)
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

	appPath := args[0]
	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}
	app, err := application.Parse("massdriver.yaml")
	if err != nil {
		return err
	}
	if errType := checkIsApplication(app); errType != nil {
		return errType
	}

	if errPub := application.Publish(c, appPath); errPub != nil {
		return errPub
	}
	log.Info().Msg("Application published successfully!")

	return nil
}

// TODO: move to "mass repo"
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
