package cmd

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const configFile = "massdriver.yaml"

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

var bundleGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generates a new bundle",
	RunE:    runBundleGenerate,
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

	bundleCmd.AddCommand(bundleGenerateCmd)
	bundleGenerateCmd.Flags().StringP("output-dir", "o", ".", "Directory to generate bundle in")

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
		output = filepath.Dir(configFile)
	}

	log.Info().Msg("building bundle")
	bun, err := bundle.Parse(configFile, nil)
	if errBuild := bun.Build(c, output); errBuild != nil {
		return errBuild
	}
	log.Info().Str("output", output).Msg("bundle built")

	return err
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
		Type:         "bundle",
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

	c, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	overrides, err := getPublishOverrides(cmd)
	if err != nil {
		return err
	}

	b, err := bundle.Parse(configFile, overrides)
	if err != nil {
		return err
	}

	err = b.Hydrate(configFile, c)
	if err != nil {
		return err
	}

	err = b.GenerateSchemas(path.Dir(configFile))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = bundle.PackageBundle(configFile, &buf)
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

	fmt.Println("Bundle published successfully!")

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
