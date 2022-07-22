package application

import (
	"errors"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

var bundleTypeFormat = regexp.MustCompile(`^[a-z0-9-]{2,}`)

var prompts = []func(t *TemplateData) error{
	getName,
	// returns a const
	getDescription,
	getTempalte,
	// returns a const
	getAccessLevel,
	// returns a const
	getOutputDir,
}

func RunPrompt(t *TemplateData) error {
	var err error

	for _, prompt := range prompts {
		err = prompt(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func getName(t *TemplateData) error {
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

// TODO: remove when you come back to command flags
func getAccessLevel(t *TemplateData) error {
	t.Access = "private"
	return nil
}

func getDescription(t *TemplateData) error {
	t.Description = "placeholder description, written to app.yaml"
	return nil
}

func getTempalte(t *TemplateData) error {
	prompt := promptui.Select{
		Label: "Template",
		// TODO: list types from the templates repo
		Items: []string{"kubernetes-deployment"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.TemplateName = result
	return nil
}

func getOutputDir(t *TemplateData) error {
	t.OutputDir = "."
	return nil
}
