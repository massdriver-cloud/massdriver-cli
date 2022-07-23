/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "mass",
	Short:        "massdriver-cli",
	Long:         `Massdriver is...`,
	SilenceUsage: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// TODO: move to version file?
const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of massdriver-cli",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: check for newer versions -> Jake
		// TODO: offer to update -> Jake
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Info().Msg("massdriver-cli version " + version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "Massdriver API key (can also be set via MASSDRIVER_API_KEY environment variable)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}

func setupLogging(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
		PartsExclude: []string{
			"time",
			"level",
		},
	})
}

func initClient(cmd *cobra.Command) (*client.MassdriverClient, error) {
	c := client.NewClient()
	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return c, err
	}
	if apiKey != "" {
		c.WithAPIKey(apiKey)
	}
	return c, nil
}
