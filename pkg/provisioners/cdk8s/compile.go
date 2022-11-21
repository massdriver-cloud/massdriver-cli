package cdk8s

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
	"github.com/rs/zerolog/log"
)

// Compile a JSON Schema to Terraform Variable Definition JSON
func Compile(path string, out io.Writer) error {
	vars, varErr := getParams(path)
	if varErr != nil {
		return varErr
	}

	// You can't have an empty variable block, so if there are no vars return an empty json block
	if len(vars) == 0 {
		if _, err := out.Write([]byte("{}")); err != nil {
			return err
		}
		return nil
	}

	// variableFile := TFVariableFile{Variable: vars}
	variableFile := vars

	bytes, err := json.MarshalIndent(variableFile, "", "    ")
	if err != nil {
		return err
	}

	_, err = out.Write(append(bytes, []byte("\n")...))

	return err
}

// Compile a JSON Schema to JSON Values based on
func CompileDevParams(path string, out io.Writer) error {
	params, paramsErr := getDevParams(path)
	if paramsErr != nil {
		return fmt.Errorf("error getting dev params: %w", paramsErr)
	}

	// You can't have an empty variable block, so if there are no vars return an empty json block
	if len(params) == 0 {
		if _, err := out.Write([]byte("{}")); err != nil {
			return err
		}
		return nil
	}

	bytes, err := json.MarshalIndent(params, "", "    ")
	if err != nil {
		return err
	}

	_, err = out.Write(append(bytes, []byte("\n")...))

	return err
}

func getParams(path string) (map[string]TFVariable, error) {
	variables := map[string]TFVariable{}
	schema, err := jsonschema.GetJSONSchema(path)
	if err != nil {
		return variables, err
	}

	required := schema.Required

	for name, prop := range schema.Properties {
		variables[name] = NewTFVariable(prop, isRequired(name, required))
	}
	return variables, nil
}

func getExistingParams(path string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	abs, err := filepath.Abs(path)
	if err != nil {
		return params, err
	}

	stat, statErr := os.Stat(abs)
	if os.IsNotExist(statErr) {
		// no existing params return empty map
		return params, nil
	} else if statErr != nil {
		return params, statErr
	}
	// if the file exists but is empty
	if stat.Size() == 0 {
		return params, nil
	}

	log.Debug().Str("path", abs).Msg("reading existing params")
	byteData, err := ioutil.ReadFile(abs)
	if err != nil {
		return params, err
	}
	log.Debug().Msgf("byteData: %s", string(byteData))
	marhsalErr := json.Unmarshal(byteData, &params)
	return params, marhsalErr
}

func getDevParams(path string) (map[string]interface{}, error) {
	params, err := getExistingParams(path)
	if err != nil {
		return params, fmt.Errorf("error getting existing params: %w", err)
	}

	bundleName := filepath.Base(filepath.Dir(filepath.Dir(path)))
	namePrefix := fmt.Sprintf("local-dev-%s-000", bundleName)
	boilerplateMetadata := map[string]interface{}{
		"name_prefix": namePrefix,
		"default_tags": map[string]interface{}{
			"md-project":  "local",
			"md-target":   "dev",
			"md-manifest": bundleName,
			"md-package":  namePrefix,
		},
		"deployment": map[string]interface{}{
			"id": "local-dev-id",
		},
		"observability": map[string]interface{}{
			"alarm_webhook_url": "https://placeholder.com",
		},
	}

	// if md_metadata is not set, initialize it to a reasonable starting point
	if _, ok := params["md_metadata"]; !ok {
		// TODO name this something better than foo (e.g. the bundle name)
		params["md_metadata"] = boilerplateMetadata
	} else {
		// merge md metadata ties go to existing values
		for k, v := range boilerplateMetadata {
			if _, ok2 := params["md_metadata"].(map[string]interface{})[k]; !ok2 {
				params["md_metadata"].(map[string]interface{})[k] = v
			}
		}
	}

	// look in parent dir of schema (path for devParams will be in src/ or some bundle step dir)
	schemaPath := filepath.Join(filepath.Dir(filepath.Dir(path)), common.ParamsSchemaFilename)
	schema, err := jsonschema.GetJSONSchema(schemaPath)
	if err != nil {
		return params, fmt.Errorf("error getting schema: %w", err)
	}
	var devExample jsonschema.Example
	for _, example := range schema.Examples {
		if example.Name == "Development" {
			devExample = example
		}
	}

	// loop over top level properties
	for name, prop := range schema.Properties {
		params[name] = FillDevParam(prop, params[name], devExample.Values[name])
	}
	return params, nil
}

var placeholderValue = "TODO: REPLACE ME"

// FillDevParam fills a parameter with a development value
// this function folows the following priority for filling in values:
// 1. If the parameter is already set, use that value
// 2. If there is a 'Development' example value, use that value
// 3. If the param defines a default use that value.
// 4. If the param is an array fallback to empty array.
// 5. If the param is a number and defines a minimum use that value.
// 4. Use a TODO string placeholder value
func FillDevParam(prop jsonschema.Property, existingVal, exampleVal interface{}) interface{} { // nolint:gocognit
	// the base case is we fall back to a placeholder to indicate to the developer they should replace this value.
	var ret interface{} = placeholderValue

	// handle nested objects recursively
	if prop.Type == jsonschema.Object {
		obj := make(map[string]interface{})
		for name, nestedProp := range prop.Properties {
			valuesMap, ok := exampleVal.(map[string]interface{})
			nestedExampleValues := valuesMap[name]
			if !ok {
				if nestedProp.Type == jsonschema.Object {
					obj[name] = make(map[string]interface{})
				}
			}
			existingMap, ok := existingVal.(map[string]interface{})
			if ok {
				nestedExistingVal := existingMap[name]
				obj[name] = FillDevParam(nestedProp, nestedExistingVal, nestedExampleValues)
			} else {
				obj[name] = FillDevParam(nestedProp, nil, nestedExampleValues)
			}
		}
		return obj
	}

	if existingVal != nil {
		return existingVal
	}

	if exampleVal != nil {
		return exampleVal
	}

	if prop.Default != nil {
		return prop.Default
	}

	// fall back to an empty array
	if prop.Type == "array" {
		return []interface{}{}
	}

	if (prop.Type == "number" || prop.Type == "integer") && prop.Minimum != nil {
		return prop.Minimum
	}

	return ret
}

func isRequired(name string, required []string) bool {
	for _, elem := range required {
		if name == elem {
			return true
		}
	}
	return false
}
