package application

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

var bundleTypeFormat = regexp.MustCompile(`^[a-z0-9-]{2,}`)

var prompts = []func(t *TemplateData) error{
	getName,
	getAccessLevel,
	getDescription,
	getChart,
	getLocation,
}

func RunPrompt(t *TemplateData) error {
	var err error
	fmt.Println("in run prompt")

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

func getAccessLevel(t *TemplateData) error {
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

func getDescription(t *TemplateData) error {
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

func getChart(t *TemplateData) error {
	prompt := promptui.Select{
		Label: "Type",
		Items: []string{"application", "adhoc-job", "scheduled-job"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Chart = result
	return nil
}

func getLocation(t *TemplateData) error {
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
