/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/application"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// applicationCmd represents the application command
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Application development tools",
	Long:  ``,
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

	applicationCmd.AddCommand(applicationPublishCmd)
}

func runApplicationPublish(cmd *cobra.Command, args []string) error {
	var err error
	appPath := args[0]

	c := client.NewClient()

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithApiKey(apiKey)
	}

	app, err := application.Parse(appPath)
	if err != nil {
		return err
	}

	// Create a temporary working directory
	workingDir, err := os.MkdirTemp("", app.Name)
	if err != nil {
		return (err)
	}
	defer os.RemoveAll(workingDir)

	// Write app.yaml
	appYaml, err := os.Create(path.Join(workingDir, "app.yaml"))
	if err != nil {
		return err
	}
	defer appYaml.Close()
	appYamlBytes, err := yaml.Marshal(*app)
	if err != nil {
		return err
	}
	appYaml.Write(appYamlBytes)

	// Write bundle.yaml
	b := app.ConvertToBundle()
	bundlePath := path.Join(workingDir, "bundle.yaml")
	bundleYaml, err := os.Create(bundlePath)
	if err != nil {
		return err
	}
	defer bundleYaml.Close()
	bundleYamlBytes, err := yaml.Marshal(*b)
	if err != nil {
		return err
	}
	bundleYaml.Write(bundleYamlBytes)

	// Make src directory
	err = os.MkdirAll(path.Join(workingDir, "src"), 0744)
	if err != nil {
		return err
	}

	err = b.Hydrate(bundlePath, c)
	if err != nil {
		return err
	}

	err = b.GenerateSchemas(workingDir)
	if err != nil {
		return err
	}

	for _, step := range b.Steps {
		switch step.Provisioner {
		case "terraform":
			err = terraform.GenerateFiles(workingDir, step.Path)
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

	uploadURL, err := b.Publish(c)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = bundle.TarGzipBundle(bundlePath, &buf)
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
