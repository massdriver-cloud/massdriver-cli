package terraform

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"
)

func GenerateFiles(bundlePath string, srcDir string) error {
	massdriverVariables := map[string]interface{}{
		"variable": map[string]interface{}{
			"md_metadata": map[string]string{
				"type": "any",
			},
		},
	}

	paramsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_params_variables.tf.json"))
	if err != nil {
		return err
	}
	err = Compile(path.Join(bundlePath, common.ParamsSchemaFilename), paramsVariablesFile)
	if err != nil {
		return err
	}
	devParamsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_params.auto.tfvars.json"))
	if err != nil {
		return err
	}
	err = CompileDevParams(path.Join(bundlePath, common.ParamsSchemaFilename), devParamsVariablesFile)
	if err != nil {
		return err
	}

	connectionsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_connections_variables.tf.json"))
	if err != nil {
		return err
	}
	err = Compile(path.Join(bundlePath, common.ConnectionsSchemaFilename), connectionsVariablesFile)
	if err != nil {
		return err
	}

	massdriverVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_md_variables.tf.json"))
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(massdriverVariables, "", "    ")
	if err != nil {
		return err
	}
	_, err = massdriverVariablesFile.Write(append(bytes, []byte("\n")...))
	if err != nil {
		return err
	}

	return nil
}

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

	variableFile := TFVariableFile{Variable: vars}

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
		return paramsErr
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

func getDevParams(path string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	schema, err := jsonschema.GetJSONSchema(path)
	if err != nil {
		return params, err
	}
	var devExample jsonschema.Example
	for _, example := range schema.Examples {
		if example.Name == "Development" {
			devExample = example
		}
	}

	// loop over top level properties
	for name, prop := range schema.Properties {
		params[name] = FillDevParam(prop, devExample.Values[name])
	}
	return params, nil
}

var placeholderValue = "TODO: REPLACE ME"

// FillDevParam fills a parameter with a development value
func FillDevParam(prop jsonschema.Property, value interface{}) interface{} {
	// the base case is we fall back to a placeholder to indicate to the developer they should replace this value.
	var ret interface{} = placeholderValue

	// handle nested objects recursively
	if prop.Type == jsonschema.Object {
		obj := make(map[string]interface{})
		for name, nestedProp := range prop.Properties {
			valuesMap, ok := value.(map[string]interface{})
			nestedValues := valuesMap[name]
			if !ok {
				if nestedProp.Type == jsonschema.Object {
					obj[name] = make(map[string]interface{})
				}
			}
			obj[name] = FillDevParam(nestedProp, nestedValues)
		}
		return obj
	}

	if value != nil {
		return value
	}

	if prop.Default != nil {
		return prop.Default
	}

	// fall bactk to an empty array
	if prop.Type == "array" {
		return []interface{}{}
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
