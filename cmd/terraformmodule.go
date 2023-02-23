package cmd

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
	"github.com/massdriver-cloud/massdriver-cli/pkg/terraformmodule"
	"github.com/spf13/cobra"
)

var terraformModuleCmd = &cobra.Command{
	Use:     "terraform-module",
	Aliases: []string{"tfm"},
	Long:    ``,
}

var terraformModuleNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new terraform-module from a template",
	RunE:  runTerraformModuleNew,
}

func init() {
	rootCmd.AddCommand(terraformModuleCmd)

	terraformModuleCmd.AddCommand(terraformModuleNewCmd)
}

func runTerraformModuleNew(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templateData := template.Data{
		Access: "private",
		// TODO: unify bundle build and app build outputDir logic and support
		OutputDir: ".",
	}

	// TODO: prompt for name of module, etc
	// err := terraformmodule.RunPromptNew(&templateData)
	// if err != nil {
	// 	return err
	// }

	errGen := terraformmodule.Generate(&templateData)
	if errGen != nil {
		return errGen
	}

	return nil
}
