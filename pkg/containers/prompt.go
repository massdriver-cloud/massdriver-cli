package containers

import (
	"github.com/manifoldco/promptui"
	"github.com/massdriver-cloud/massdriver-cli/pkg/template"
)

var promptsNew = []func(t *template.Data) error{
	getName,
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
		return nil
	}

	defaultValue := "bundle"
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
