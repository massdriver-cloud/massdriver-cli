package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/config"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Bundle development tools",
	Long:  ``,
}

var bundleBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds bundle JSON Schemas",

	RunE: runBundleBuild,
}

var bundleLintCmd = &cobra.Command{
	Use:          "lint",
	Short:        "Check bundle configuration for common errors",
	SilenceUsage: true,
	RunE:         runBundleLint,
}

var bundleGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Deprecated: Generates a new bundle",
	RunE:    runBundleGenerate,
}

var bundleNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new bundle from a template",
	RunE:  runBundleGenerate,
}

var bundlePublishCmd = &cobra.Command{
	Use:          "publish",
	Short:        "Publish a bundle to Massdriver",
	RunE:         runBundlePublish,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", "", "Path to output directory (default is massdriver.yaml directory)")

	bundleCmd.AddCommand(bundleLintCmd)

	bundleCmd.AddCommand(bundleGenerateCmd)
	bundleGenerateCmd.Flags().StringP("output-dir", "o", ".", "Directory to generate bundle in")
	bundleCmd.AddCommand(bundleNewCmd)
	bundleNewCmd.Flags().StringP("output-dir", "o", ".", "Directory to generate bundle in")

	bundleCmd.AddCommand(bundlePublishCmd)
	bundlePublishCmd.Flags().String("access", "", "Override the access, useful in CI for deploying to sandboxes.")
}

func runBundleBuild(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	// default the output to the path of the massdriver.yaml file
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	if output == "" {
		output = filepath.Dir(common.MassdriverYamlFilename)
	}

	log.Info().Msg("building bundle")
	b, err := bundle.Parse(common.MassdriverYamlFilename, nil)
	if err != nil {
		return err
	}
	if !b.IsInfrastructure() {
		return fmt.Errorf("this command can only be used with bundle type 'infrastructure'")
	}

	if errBuild := b.Build(c, output); errBuild != nil {
		return errBuild
	}
	log.Info().Str("output", output).Msg("bundle built")

	return nil
}

func runBundleGenerate(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	var err error

	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return err
	}

	templateData := template.Data{
		OutputDir:    outputDir,
		Type:         "infrastructure",
		TemplateName: "terraform",
	}

	err = bundle.RunPrompt(&templateData)
	if err != nil {
		return err
	}

	err = bundle.Generate(&templateData)
	if err != nil {
		return err
	}

	return nil
}

func runBundlePublish(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	conf := config.Get()

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	overrides, err := getPublishOverrides(cmd)
	if err != nil {
		return err
	}
	b, err := bundle.Parse(common.MassdriverYamlFilename, overrides)
	if err != nil {
		return err
	}
	if !b.IsInfrastructure() {
		return fmt.Errorf("this command can only be used with bundle type 'infrastructure'")
	}

	b.Hydrate(common.MassdriverYamlFilename, c)
	if err != nil {
		return err
	}

	if errLint := bundle.Lint(b); errLint != nil {
		msg := fmt.Sprintf("Warning! Bundle linting failed %s\nForce flag enabled, proceeding with publish", errLint.Error())
		log.Warn().Msg(msg)
	}

	if errPublish := bundle.Publish(c, b); errPublish != nil {
		return errPublish
	}

	msg := fmt.Sprintf("%s %s '%s' published successfully", b.Access, b.Type, b.Name)
	log.Info().Str("organizationId", conf.OrgID).Msg(msg)
	return nil
}

// runBundleLint checks the bundle for common configuration errors
func runBundleLint(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	b, err := bundle.Parse(common.MassdriverYamlFilename, nil)
	if err != nil {
		return err
	}
	if !b.IsInfrastructure() {
		return fmt.Errorf("this command can only be used with bundle type 'infrastructure'")
	}

	b.Hydrate(common.MassdriverYamlFilename, c)
	if err != nil {
		return err
	}

	err = bundle.Lint(b)
	if err != nil {
		return err
	}

	log.Info().Msg("Configuration is valid!")

	return nil
}

func getPublishOverrides(cmd *cobra.Command) (map[string]interface{}, error) {
	access, err := cmd.Flags().GetString("access")
	if err != nil {
		return nil, err
	}

	overrides := map[string]interface{}{
		"access": access,
	}

	return overrides, nil
}
