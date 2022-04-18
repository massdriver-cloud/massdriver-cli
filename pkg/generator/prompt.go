package generator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

var bundleTypeFormat = regexp.MustCompile(`^[a-z0-9-]{5,}`)

var prompts = []func(t *TemplateData) error{
	getName,
	getType,
	getAccessLevel,
	getDescription,
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

func getType(t *TemplateData) error {
	validate := func(input string) error {
		if !bundleTypeFormat.MatchString(input) {
			return errors.New("name must be greater than 4 characters and can only include lowercase letters and dashes")
		}
		return nil
	}

	defaultValue := strings.Replace(strings.ToLower(t.Name), " ", "-", -1)

	prompt := promptui.Prompt{
		Label:    "Type",
		Validate: validate,
		Default:  defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	t.Type = result
	return nil
}

func getName(t *TemplateData) error {
	prompt := promptui.Prompt{
		Label: "Name",
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
