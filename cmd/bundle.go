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

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Bundle development tools",
	Long:  ``,
}

var bundleBuildCmd = &cobra.Command{
	Use:   "build [Path to massdriver.yaml]",
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
	Use:          "publish [Path to massdriver.yaml]",
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
}

func runBundleBuild(cmd *cobra.Command, args []string) error {
	var err error
	var bundlePath string

	if len(args) == 0 {
		bundlePath = "massdriver.yaml"
	} else {
		bundlePath = args[0]
	}

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
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while building bundle")
		return err
	}
	if output == "" {
		output = filepath.Dir(bundlePath)
	}

	log.Info().Str("bundle", bundlePath).Msg("building bundle")

	b, err := bundle.Parse(bundlePath)
	if err != nil {
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while parsing bundle")
		return err
	}

	err = b.Hydrate(bundlePath, c)
	if err != nil {
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while hydrating bundle")
		return err
	}

	err = b.GenerateSchemas(output)
	if err != nil {
		log.Error().Err(err).Str("bundle", bundlePath).Msg("an error occurred while generating bundle schema files")
		return err
	}

	for _, step := range b.Steps {
		switch step.Provisioner {
		case "terraform":
			err = terraform.GenerateFiles(output, step.Path)
			if err != nil {
				log.Error().Err(err).Str("bundle", bundlePath).Str("provisioner", step.Provisioner).Msg("an error occurred while generating provisioner files")
				return err
			}
		case "exec":
			// No-op (Golang doesn't not fallthrough unless explicitly stated)
		default:
			log.Error().Str("bundle", bundlePath).Msg("unknown provisioner: " + step.Provisioner)
			return fmt.Errorf("unknown provisioner: %v", step.Provisioner)
		}
	}

	log.Info().Str("bundle", bundlePath).Str("output", output).Msg("bundle built")

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
	var bundlePath string

	if len(args) == 0 {
		bundlePath = "massdriver.yaml"
	} else {
		bundlePath = args[0]
	}

	c := client.NewClient()

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithApiKey(apiKey)
	}

	b, err := bundle.Parse(bundlePath)
	if err != nil {
		return err
	}

	err = b.Hydrate(bundlePath, c)
	if err != nil {
		return err
	}

	err = b.GenerateSchemas(path.Dir(bundlePath))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = bundle.PackageBundle(bundlePath, &buf)
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
