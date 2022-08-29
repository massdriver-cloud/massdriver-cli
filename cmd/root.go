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
	Use:   "mass",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "Massdriver API key (can also be set via MASSDRIVER_API_KEY environment variable)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("json", false, "Enable JSON output")
}

func setupLogging(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	json, _ := cmd.Flags().GetBool("json")

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if !json {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stdout,
			PartsExclude: []string{
				"time",
			},
		})
	}
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
