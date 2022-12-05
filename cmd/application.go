package cmd

import (
	"fmt"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/config"
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
	Use:   "build",
	Short: "Builds the app for local development",
	RunE:  runApplicationBuild,
}

var applicationNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new application from a template",
	RunE:  runApplicationNew,
}

var applicationPublishCmd = &cobra.Command{
	Use:          "publish",
	Short:        "Publish an application to Massdriver",
	RunE:         runApplicationPublish,
	SilenceUsage: true,
}

var applicationDeployCmd = &cobra.Command{
	Use:   `deploy <project>-<target>-<manifest>`,
	Short: "Deploy an application on Massdriver",
	Long:  `Deploy an application package on Massdriver. This application must be published to Massdriver first and a package configured in a particular target. This command requires setting MASSDRIVER_ORG_ID environment variable. For more info see https://docs.massdriver.cloud/applications/deploying`,
	Args:  cobra.ExactArgs(1),
	RunE:  RunApplicationDeploy,
}

var applicationTemplatesCmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"tmpl"},
	Short:   "Application template development tools",
	Long:    ``,
}

var applicationTemplatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists available application templates",
	RunE:  runApplicationTemplatesList,
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
	applicationCmd.AddCommand(applicationDeployCmd)
	applicationCmd.AddCommand(applicationTemplatesCmd)

	applicationTemplatesCmd.AddCommand(templatesRefreshCmd)
	applicationTemplatesCmd.AddCommand(applicationTemplatesListCmd)
}

func runApplicationBuild(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}
	// TODO: app/bundle build directories
	output := "."

	app, err := application.Parse(common.MassdriverYamlFilename, nil)
	if err != nil {
		return err
	}
	if !app.IsApplication() {
		return fmt.Errorf("this command can only be used with bundle type 'application'")
	}
	application.Build(app, c, output)

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

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}
	app, err := application.Parse(common.MassdriverYamlFilename, nil)
	if err != nil {
		return err
	}
	if !app.IsApplication() {
		return fmt.Errorf("this command can only be used with bundle type 'application'")
	}

	if errPub := bundle.Publish(c, app); errPub != nil {
		return errPub
	}
	log.Info().Msg("Application published successfully!")

	return nil
}

// RunApplicationDeploy deploys an application to Massdriver (exported to avoid code duplication for deprecated `mass deploy` command)
func RunApplicationDeploy(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	app, err := application.Parse(common.MassdriverYamlFilename, nil)
	if err != nil {
		return err
	}
	if !app.IsApplication() {
		return fmt.Errorf("this command can only be used with bundle type 'application'")
	}
	// original
	name := args[0]

	c := config.Get()

	client := api.NewClient()
	client2 := api2.NewClient(c.APIKey)
	deployment, err := api.DeployPackage(client, &client2, c.OrgID, name)

	if err != nil {
		log.Fatal().Err(err).Str("deploymentId", deployment.ID).Msg("deployment failed")
		return err
	}

	log.Info().Str("deploymentId", deployment.ID).Msgf("Deployment %s", strings.ToLower(deployment.Status))
	return nil
}

// TODO: move to `mass repos`
func runApplicationTemplatesList(cmd *cobra.Command, args []string) error {
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
