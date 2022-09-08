package application

import (
	"errors"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/cache"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

var bundleTypeFormat = regexp.MustCompile(`^[a-z0-9-]{2,}`)

var promptsNew = []func(t *template.Data) error{
	getName,
	getDescription,
	getAccessLevel,
	getTemplate,
	getOutputDir,
	bundle.GetConnections,
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

var IgnoredTemplateDirs = []string{"alpha"}

func getTemplate(t *template.Data) error {
	templates, err := cache.ApplicationTemplates()
	templates = removeIgnoredTemplateDirectories(templates, IgnoredTemplateDirs)

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

func removeIgnoredTemplateDirectories(templates []string, remove []string) []string {
	for i := 0; i < len(templates); i++ {
		template := templates[i]
		for _, rem := range remove {
			if template == rem {
				templates = append(templates[:i], templates[i+1:]...)
				i--
				break
			}
		}
	}

	return templates
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
