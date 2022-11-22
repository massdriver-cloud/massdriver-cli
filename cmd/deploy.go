/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// TODO remove mass deploy command entirely during next major release.
var deployCmd = &cobra.Command{
	Use:        "deploy",
	Short:      "Deploy a configured package",
	Long:       ``,
	Args:       cobra.ExactArgs(1),
	RunE:       RunApplicationDeploy,
	Deprecated: "`mass deploy` is deprecated and will be removed in a future release; use `mass app deploy` instead",
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
