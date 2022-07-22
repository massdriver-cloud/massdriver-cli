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
	getTempalte,
	// returns a const, leaving for now
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
	if t.Name != "" {
		return nil
	}

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

func getTempalte(t *TemplateData) error {
	if t.TemplateName != "" {
		return nil
	}

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
