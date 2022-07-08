/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:     "application",
	Aliases: []string{"app"},
	Short:   "Application development tools",
	Long:    ``,
}

var applicationGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generates a new application template",
	RunE:    runApplicationGenerate,
}

var applicationPublishCmd = &cobra.Command{
	Use:          "publish [Path to app.yaml]",
	Short:        "Publish an application to Massdriver",
	Args:         cobra.ExactArgs(1),
	RunE:         runApplicationPublish,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(applicationCmd)

	applicationCmd.AddCommand(applicationGenerateCmd)
	applicationCmd.AddCommand(applicationPublishCmd)
}

func runApplicationGenerate(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templateData := application.TemplateData{}

	err := application.RunPrompt(&templateData)
	if err != nil {
		return err
	}

	err = application.Generate(&templateData)
	if err != nil {
		return err
	}

	return nil
}

func runApplicationPublish(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	var err error
	appPath := args[0]

	c := client.NewClient()

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithAPIKey(apiKey)
	}

	// Create a temporary working directory
	workingDir, err := os.MkdirTemp("", "application")
	if err != nil {
		return (err)
	}
	defer os.RemoveAll(workingDir)

	var buf bytes.Buffer
	b, err := application.PackageApplication(appPath, c, workingDir, &buf)
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

	fmt.Println("Application published successfully!")

	return nil
}
