package application

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/itchyny/gojq"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/schema2json"
)

func Lint(app *bundle.Bundle) error {
	err := bundle.Lint(app)
	if err != nil {
		return err
	}

	_, err = LintEnvs(app)
	if err != nil {
		return fmt.Errorf("an error occurred while validating the envs: %s", err.Error())
	}

	return nil
}

func LintEnvs(app *bundle.Bundle) (map[string]string, error) {
	result := map[string]string{}

	input, err := buildEnvsInput(app)
	if err != nil {
		return nil, fmt.Errorf("error building env query: %s", err.Error())
	}

	for name, query := range app.App.Envs {
		jq, err := gojq.Parse(query)
		if err != nil {
			return result, errors.New("The jq query for environment variable " + name + " is invalid: " + err.Error())
		}

		iter := jq.Run(input)
		value, ok := iter.Next()
		if !ok || value == nil {
			return result, errors.New("The jq query for environment variable " + name + " didn't produce a result")
		}
		if err, ok := value.(error); ok {
			return result, errors.New("The jq query for environment variable " + name + " produced an error: " + err.Error())
		}
		var valueString string
		if valueString, ok = value.(string); !ok {
			resultBytes, err := json.Marshal(value)
			if err != nil {
				return result, errors.New("The jq query for environment variable " + name + " produced an uninterpretable value: " + err.Error())
			}
			valueString = string(resultBytes)
		}
		_, multiple := iter.Next()
		if multiple {
			return result, errors.New("The jq query for environment variable " + name + " produced multiple values, which isn't supported")
		}
		result[name] = valueString
	}

	return result, nil
}

func buildEnvsInput(b *bundle.Bundle) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	paramsSchema, err := schema2json.ParseMapStringInterface(b.Params)
	if err != nil {
		return nil, err
	}
	connectionsSchema, err := schema2json.ParseMapStringInterface(b.Connections)
	if err != nil {
		return nil, err
	}
	result["params"], err = schema2json.GenerateJSON(paramsSchema)
	if err != nil {
		return nil, err
	}
	result["connections"], err = schema2json.GenerateJSON(connectionsSchema)
	if err != nil {
		return nil, err
	}

	secrets := map[string]interface{}{}
	for name := range b.App.Secrets {
		secrets[name] = "some-secret-value"
	}
	result["secrets"] = secrets

	// by, _ := json.Marshal(b.Connections)
	// fmt.Println(string(by))

	return result, nil
}
