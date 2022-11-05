package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
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

	bundleCmd.AddCommand(bundleGenerateCmd)
	bundleGenerateCmd.Flags().StringP("output", "o", ".", "Directory to generate bundle in")
	bundleCmd.AddCommand(bundleNewCmd)
	// --name ${bundleName} --access ${bundleAccess} --description ${bundleDescription} --output ${bundleOutput} --dependencies
	bundleNewCmd.Flags().StringP("name", "", ".", "Name of the bundle")
	bundleNewCmd.Flags().StringP("access", "", ".", "")
	bundleNewCmd.Flags().StringP("description", "", ".", "")
	bundleNewCmd.Flags().StringP("connections", "", ".", "")
	bundleNewCmd.Flags().StringP("output", "", ".", "Directory to generate bundle in")

	bundleCmd.AddCommand(bundlePublishCmd)
	bundlePublishCmd.Flags().String("access", "", "Override the access, useful in CI for deploying to sandboxes.")
}

func checkIsBundle(b *bundle.Bundle) error {
	if b.Type != "bundle" {
		return fmt.Errorf("mass bundle publish can only be used with bundle type massdriver.yaml")
	}
	return nil
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
	if errType := checkIsBundle(b); errType != nil {
		return errType
	}

	if errBuild := b.Build(c, output); errBuild != nil {
		return errBuild
	}
	log.Info().Str("output", output).Msg("bundle built")

	return nil
}

func runBundleGenerate(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	name, errName := cmd.Flags().GetString("name")
	if errName != nil {
		return errName
	}

	access, errAcc := cmd.Flags().GetString("access")
	if errAcc != nil {
		return errAcc
	}

	description, errDesc := cmd.Flags().GetString("description")
	if errDesc != nil {
		return errDesc
	}

	// connections, errCon := cmd.Flags().GetString("connections")
	// if errCon != nil {
	// 	return errCon
	// }

	outputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	templateData := template.Data{
		Name: name,
		Access: access,
		Description: description,
		Connections: map[string]interface{}{
			"hank": "hank",
		},
		OutputDir:    outputDir,
		Type:         "bundle",
		TemplateName: "terraform",
	}
	fmt.Printf("Tempalte is %v", templateData)

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
	b, err := bundle.Parse(common.MassdriverYamlFilename, overrides)
	if err != nil {
		return err
	}
	if errType := checkIsBundle(b); errType != nil {
		return errType
	}

	if errPublish := bundle.Publish(c, b); errPublish != nil {
		return errPublish
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
