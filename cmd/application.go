/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var applicationPublishCmd = &cobra.Command{
	Use:                   "publish",
	Short:                 "Publish application to Massdriver",
	RunE:                  RunApplicationCreate,
	DisableFlagsInUseLine: true,
}

var applicationLintCmd = &cobra.Command{
	Use:                   "lint",
	Short:                 "Lint the local application definition to ensure proper formatting",
	RunE:                  RunApplicationValidate,
	DisableFlagsInUseLine: true,
}

var applicationParseCmd = &cobra.Command{
	Use:                   "parse [Path to app.yaml]",
	Short:                 "Parses and app.yaml file",
	Args:                  cobra.ExactArgs(1),
	RunE:                  RunApplicationParse,
	DisableFlagsInUseLine: true,
}

var generateCmd = &cobra.Command{
	Use:                   "generate",
	Short:                 "DELETE ME",
	RunE:                  RunGenerate,
	DisableFlagsInUseLine: true,
}

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:     "application",
	Aliases: []string{"app"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	rootCmd.AddCommand(applicationCmd)

	applicationCmd.AddCommand(applicationPublishCmd)
	applicationPublishCmd.Flags().StringP("file", "f", "app.yaml", "Application config file")

	applicationCmd.AddCommand(applicationLintCmd)
	applicationLintCmd.Flags().StringP("file", "f", "app.yaml", "Application config file")
	applicationLintCmd.Flags().StringP("schema", "s", "", "Schema file")

	applicationCmd.AddCommand(applicationParseCmd)
}

func RunApplicationCreate(cmd *cobra.Command, args []string) error {
	return nil
}

func RunApplicationValidate(cmd *cobra.Command, args []string) error {
	config, _ := cmd.Flags().GetString("file")
	schema, _ := cmd.Flags().GetString("file")

	valid, err := application.Lint(config, schema)
	if err != nil {
		return err
	}
	if !valid {
		os.Exit(1)
	}
	return nil
}

func RunApplicationParse(cmd *cobra.Command, args []string) error {
	path := args[0]

	log.Info().Msg("Parsing application")

	application.Parse(path)

	return nil
}

func RunGenerate(cmd *cobra.Command, args []string) error {
	application.Generate()
	return nil
}
