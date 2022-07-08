/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/definition"
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var definitionCmd = &cobra.Command{
	Use:   "definition",
	Short: "artifact definition management",
	Long:  ``,
}

var definitionGetCmd = &cobra.Command{
	Use:          "get [definition]",
	Short:        "Get an artifact definition from Massdriver",
	Args:         cobra.ExactArgs(1),
	RunE:         runDefinitionGet,
	SilenceUsage: true,
}

var definitionPublishCmd = &cobra.Command{
	Use:          "publish",
	Short:        "Publish an artifact definition to Massdriver",
	RunE:         runDefinitionPublish,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(definitionCmd)

	definitionCmd.AddCommand(definitionGetCmd)

	definitionCmd.AddCommand(definitionPublishCmd)
	definitionPublishCmd.Flags().StringP("file", "f", "", "File containing artifact definition schema (use - for stdin)")
	if err := definitionPublishCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
}

func runDefinitionGet(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c := client.NewClient()

	defName := args[0]

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithAPIKey(apiKey)
	}

	def, err := definition.GetDefinition(c, defName)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(def)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}

func runDefinitionPublish(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	c := client.NewClient()

	defPath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey != "" {
		c.WithAPIKey(apiKey)
	}

	var defFile *os.File
	if defPath == "-" {
		defFile = os.Stdin
	} else {
		defFile, err = os.Open(defPath)
		if err != nil {
			fmt.Println(err)
		}
		defer defFile.Close()
	}

	byteValue, _ := ioutil.ReadAll(defFile)
	var art definition.Definition

	if jsonErr := json.Unmarshal(byteValue, &art); err != nil {
		return jsonErr
	}

	if pubErr := art.Publish(c); pubErr != nil {
		return pubErr
	}

	fmt.Println("Definition published successfully!")

	return nil
}
