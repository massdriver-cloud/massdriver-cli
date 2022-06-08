package cmd

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"

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
	var err error

	c := client.NewClient()

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithApiKey(apiKey)
	}

	// default the output to the path of the massdriver.yaml file
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while building bundle")
		return err
	}
	if output == "" {
		output = filepath.Dir(configFile)
	}

	log.Info().Msg("building bundle")

	b, err := bundle.Parse(configFile, nil)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while parsing bundle")
		return err
	}

	err = b.Hydrate(configFile, c)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while hydrating bundle")
		return err
	}

	err = b.GenerateSchemas(output)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while generating bundle schema files")
		return err
	}

	for _, step := range b.Steps {
		switch step.Provisioner {
		case "terraform":
			err = terraform.GenerateFiles(output, step.Path)
			if err != nil {
				log.Error().Err(err).Str("provisioner", step.Provisioner).Msg("an error occurred while generating provisioner files")
				return err
			}
		default:
			log.Error().Msg("unknown provisioner: " + step.Provisioner)
			return fmt.Errorf("unknown provisioner: %v", step.Provisioner)
		}
	}

	log.Info().Str("output", output).Msg("bundle built")

	return err
}

func runBundleGenerate(cmd *cobra.Command, args []string) error {
	var err error

	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return err
	}

	templateData := bundle.TemplateData{
		OutputDir: outputDir,
		Type:      "bundle",
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
	var err error
	c := client.NewClient()

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithApiKey(apiKey)
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
