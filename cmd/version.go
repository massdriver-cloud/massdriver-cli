/*

*/
package cmd

import (
	"fmt"

	"github.com/massdriver-cloud/massdriver-cli/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "version of mass cli",
	Long:    ``,
	Run:     runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	isOld, latest, err := version.CheckForNewerVersionAvailable()
	if err != nil {
		fmt.Printf("could not check for newer versions at %v: %v. skipping...\n", version.LatestReleaseURL, err.Error())
	} else if isOld {
		fmt.Printf("A newer version of the CLI is available, you can download it here: %v\n", version.LatestReleaseURL)
	}
	fmt.Printf("mass version: %v (git SHA: %v) \n", version.MassVersion(), version.MassGitSHA())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
