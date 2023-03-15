package bundle

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"

	"github.com/manifoldco/promptui"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"
)

const noneDep = "(None)"

var bundleTypeFormat = regexp.MustCompile(`^[a-z0-9-]{5,}`)
var connectionNameFormat = regexp.MustCompile(`^[a-z]+[a-z0-9_]*[a-z0-9]+$`)

var prompts = []func(t *template.Data) error{
	getName,
	getAccessLevel,
	getDescription,
	getOutputDir,
	GetConnections,
}

func RunPrompt(t *template.Data) error {
	var err error

	for _, prompt := range prompts {
		err = prompt(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func getName(t *template.Data) error {
	validate := func(input string) error {
		if !bundleTypeFormat.MatchString(input) {
			return errors.New("name must be greater than 4 characters and can only include lowercase letters and dashes")
		}
		return nil
	}

	defaultValue := strings.ReplaceAll(strings.ToLower(t.Name), " ", "-")

	prompt := promptui.Prompt{
		Label:    "Name",
		Validate: validate,
		Default:  defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	t.Name = result
	return nil
}

func getAccessLevel(t *template.Data) error {
	prompt := promptui.Select{
		Label: "Access Level",
		Items: []string{"public", "private"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Access = result
	return nil
}

func getDescription(t *template.Data) error {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("Description cannot be empty.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Description",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Description = result
	return nil
}

func getOutputDir(t *template.Data) error {
	prompt := promptui.Prompt{
		Label:   `Output directory`,
		Default: t.Name,
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.OutputDir = result
	return nil
}

// TODO fetch these from the API instead of hardcoding
func GetConnections(t *template.Data) error {
	artDefs, err := common.ListMassdriverArtifactDefinitions()
	if err != nil {
		return err
	}
	var selectedDeps []string
	multiselect := &survey.MultiSelect{
		Message: "What connections do you need?\n  If you don't need any, just hit enter or select (None)\n",
		Options: append([]string{noneDep}, artDefs...),
	}
	err = survey.AskOne(multiselect, &selectedDeps)
	if err != nil {
		return err
	}
	var depMap []template.Connection
	for i, v := range selectedDeps {
		if v == noneDep {
			if len(selectedDeps) > 1 {
				return fmt.Errorf("if selecting %v, you cannot select other dependecies. selected %#v", noneDep, selectedDeps)
			}
			return nil
		}

		validate := func(input string) error {
			if !connectionNameFormat.MatchString(input) {
				return errors.New("name must be at least 2 characters, start with a-z, use lowercase letters, numbers and underscores. It can not end with an underscore")
			}
			return nil
		}

		log.Info().Msgf("Please enter a name for the connection: \"%v\"\nThis will be the variable name used to reference it in your app|bundle IaC", v)
		prompt := promptui.Prompt{
			Label:    `Name`,
			Validate: validate,
		}

		result, errName := prompt.Run()
		if errName != nil {
			return errName
		}

		depMap = append(depMap, template.Connection{Name: result, ArtifactDefinition: selectedDeps[i]})
	}

	t.Connections = depMap
	return nil
}
