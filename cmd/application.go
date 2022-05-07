/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
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
	applicationPublishCmd.Flags().StringP("api-key", "k", "", "Massdriver API key (can also be set via MASSDRIVER_API_KEY environment variable)")
}

func runApplicationPublish(cmd *cobra.Command, args []string) error {
	// var err error
	// applicationPath := args[0]

	// c := client.NewClient()

	// apiKey, err := cmd.Flags().GetString("api-key")
	// if err != nil {
	// 	return err
	// }
	// if apiKey != "" {
	// 	c.WithApiKey(apiKey)
	// }

	// app, err := application.Parse(applicationPath)
	// if err != nil {
	// 	return err
	// }

	// b, err := app.ConvertToBundle(c)
	// if err != nil {
	// 	return err
	// }

	// uploadURL, err := b.Publish(c)
	// if err != nil {
	// 	return err
	// }

	// var buf bytes.Buffer
	// err = application.TarGzipApplication(applicationPath, &buf)
	// if err != nil {
	// 	return err
	// }

	// err = bundle.UploadToPresignedS3URL(uploadURL, &buf)
	// if err != nil {
	// 	return err
	// }

	return nil
}
