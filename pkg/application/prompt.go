package application

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

const noneDep = "(None)"

var bundleTypeFormat = regexp.MustCompile(`^[a-z0-9-]{2,}`)

var prompts = []func(t *template.Data) error{
	getName,
	getDescription,
	getAccessLevel,
	getTemplate,
	// TODO: deprecate
	getChart,
	getLocation,
}

var promptsNew = []func(t *template.Data) error{
	getName,
	getDescription,
	getAccessLevel,
	getTemplate,
	getOutputDir,
	getDeps,
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

func RunPromptNew(t *template.Data) error {
	var err error

	for _, prompt := range promptsNew {
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
			return errors.New("name must be 2 or more characters and can only include lowercase letters and dashes")
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
	if t.Access != "" {
		return nil
	}

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
	prompt := promptui.Prompt{
		Label: "Description",
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Description = result
	return nil
}

// TODO: deprecate
func getChart(t *template.Data) error {
	prompt := promptui.Select{
		Label: "Access Level",
		Items: []string{"application", "adhoc-job", "scheduled-job"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Chart = result
	return nil
}

func getTemplate(t *template.Data) error {
	templates, err := cache.ApplicationTemplates()
	if err != nil {
		return err
	}
	prompt := promptui.Select{
		Label: "Template",
		Items: templates,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.TemplateName = result
	return nil
}

// TODO: deprecate
func getLocation(t *template.Data) error {
	prompt := promptui.Prompt{
		Label:     "Chart Location",
		Default:   "./chart",
		AllowEdit: true,
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Location = result
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

func getDeps(t *template.Data) error {
	artifacts, err := GetMassdriverArtifacts()
	if err != nil {
		return err
	}
	var selectedDeps []string
	multiselect := &survey.MultiSelect{
		Message: "What artifacts does your application depend on? If you have no dependencies just hit enter or only select (None)",
		Options: append([]string{noneDep}, artifacts...),
	}
	err = survey.AskOne(multiselect, &selectedDeps)
	if err != nil {
		return err
	}
	depMap := make(map[string]string)
	for i, v := range selectedDeps {
		if v == noneDep {
			t.Dependencies = make(map[string]string)
			if len(selectedDeps) > 1 {
				return fmt.Errorf("if selecting %v, you cannot select other dependecies. selected %#v", noneDep, selectedDeps)
			}
			return nil
		}
		// TODO may have to replace the slash in artifact names
		// dependencies are a map with indexed key so in the future we could allow selecting multiple of the same artifact type
		depMap[fmt.Sprintf("%v_%v", v, i)] = selectedDeps[i]
	}
	t.Dependencies = depMap
	return nil
}
