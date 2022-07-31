/*

 */
package cmd

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/version"
	"github.com/rs/zerolog/log"
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
	isOld, _, err := version.CheckForNewerVersionAvailable()
	if err != nil {
		log.Info().Msgf("could not check for newer versions at %v: %v. skipping...\n", version.LatestReleaseURL, err.Error())
	} else if isOld {
		log.Info().Msgf("A newer version of the CLI is available, you can download it here: %v\n", version.LatestReleaseURL)
	}
	log.Info().Msgf("mass version: %v (git SHA: %v) \n", version.MassVersion(), version.MassGitSHA())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
