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

var bundleNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new bundle from a template",
	RunE:  runBundleNew,
}

var bundlePublishCmd = &cobra.Command{
	Use:          "publish",
	Short:        "Publish a bundle to Massdriver",
	RunE:         runBundlePublish,
	SilenceUsage: true,
}

var name string
var access string
var description string
var output string
var connections []string
var connectionNames []string

func init() {
	rootCmd.AddCommand(bundleCmd)

	bundleCmd.AddCommand(bundleBuildCmd)
	bundleBuildCmd.Flags().StringP("output", "o", "", "Path to output directory (default is massdriver.yaml directory)")

	bundleCmd.AddCommand(bundleNewCmd)
	bundleNewCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the bundle")
	bundleNewCmd.Flags().StringVarP(&access, "access", "a", "", "Access level of the bundle")
	bundleNewCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the bundle")
	bundleNewCmd.Flags().StringVarP(&output, "output", "o", "", "Path to output directory")
	bundleNewCmd.Flags().StringSliceVarP(&connections, "connections", "c", []string{}, "List of bundle connections")
	bundleNewCmd.Flags().StringSliceVarP(&connectionNames, "connection-names", "x", []string{}, "Names for the bundle connections")

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

func runBundleNew(cmd *cobra.Command, args []string) error {
	// Check the value of the name flag
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get name flag")
	}
	if name == "" {
		// If the name flag is not set, prompt the user for the values of the name, access, description, output, connections, and connection names flags
		templateData := template.Data{
			Type:         "bundle",
			TemplateName: "terraform",
		}
		return bundle.RunPrompt(&templateData)
	}

	if len(name) < 5 {
		log.Fatal().Msg("name must be at least 5 characters long")
	}

	// If the name flag is set, then check the values of the access, description, output, connections, and connection names flags
	access, err := cmd.Flags().GetString("access")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get access flag")
	}
	if access == "" {
		log.Fatal().Msg("access flag is required")
	}
	if access != "private" && access != "public" {
		log.Fatal().Msg("access must be either 'private' or 'public'")
	}

	description, err := cmd.Flags().GetString("description")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get description flag")
	}
	if description == "" {
		log.Fatal().Msg("description flag is required")
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get output flag")
	}
	if output == "" {
		log.Fatal().Msg("output flag is required")
	}

	connections, err := cmd.Flags().GetStringSlice("connections")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get connections flag")
	}
	if len(connections) == 0 {
		log.Fatal().Msg("connections flag is required")
	}

	connectionNames, err := cmd.Flags().GetStringSlice("connection-names")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get connection-names flag")
	}

	// Check that the number of connections and connection names is equal so they can be indexed to match
	if len(connections) != len(connectionNames) {
		log.Fatal().Msg("the number of connections and connection names must be equal")
	}

	// Create a map of connection names to connections
	connectionMap := make(map[string]string)
	for i, connection := range connections {
		connectionMap[connectionNames[i]] = connection
	}

	// If all flags are set, generate the bundle
	templateData := template.Data{
		OutputDir:       output,
		Type:            "bundle",
		TemplateName:    "terraform",
		Name:            name,
		Access:          access,
		Description:     description,
		Connections:     connections,
		ConnectionNames: connectionNames,
	}
	return bundle.Generate(&templateData)
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
	if !b.IsInfrastructure() {
		return fmt.Errorf("this command can only be used with bundle type 'infrastructure'")
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
