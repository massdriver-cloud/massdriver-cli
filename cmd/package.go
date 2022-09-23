package cmd

import (
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "Pacakge development tools",
	Long:    ``,
}

var packageConfigureCmd = &cobra.Command{
	Use:     `configure <project>-<target>-<manifest>`,
	Aliases: []string{"config"},
	Short:   "Configure a package in a Massdriver target",
	Long:    ``,
	RunE:    runPackageConfigure,
	Args:    cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(applicationCmd)
	packageCmd.AddCommand(packageConfigureCmd)
}

func runPackageConfigure(cmd *cobra.Command, args []string) error {
	return nil
}
