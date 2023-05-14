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

// var useApplicationTemplate bool

func init() {
	rootCmd.AddCommand(terraformModuleCmd)

	terraformModuleCmd.AddCommand(terraformModuleNewCmd)
	terraformModuleNewCmd.Flags().BoolP("application", "a", false, "Is this an application Terraform module?")
}

func runTerraformModuleNew(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	templateData := template.Data{
		OutputDir: "terraform-module",
	}

	useApplicationTemplate, err := cmd.Flags().GetBool("application")
	if err != nil {
		return err
	}

	if useApplicationTemplate {
		// TODO: prompt for app things
		// err := terraformmodule.RunPromptNew(&templateData)
		// if err != nil {
		// 	return err
		// }
		templateData.TemplateName = "massdriver-application"
		templateData.OutputDir = "massdriver-application"

		errGen := terraformmodule.GenerateApplication(&templateData)
		if errGen != nil {
			return errGen
		}
		return nil
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
