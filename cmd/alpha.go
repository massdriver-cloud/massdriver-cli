package cmd

import (
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "command group for alpha features",
	Long:  "command group for alpha features of this CLI that are subject to breaking changes",
}

func init() {
	rootCmd.AddCommand(alphaCmd)
}
