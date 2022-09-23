package terraform

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

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
	devParamsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, common.DevParamsFilename))
	if err != nil {
		return err
	}
	err = CompileDevParams(path.Join(bundlePath, common.DevParamsFilename), devParamsVariablesFile)
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

func getExistingParams(path string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// no existing params return empty map
		return params, nil
	}
	byteData, err := ioutil.ReadFile(path)
	if err != nil {
		return params, err
	}
	marhsalErr := json.Unmarshal(byteData, &params)
	return params, marhsalErr
}

func getDevParams(path string) (map[string]interface{}, error) {
	params, err := getExistingParams(path)
	if err != nil {
		return params, err
	}
	schema, err := jsonschema.GetJSONSchema(strings.Replace(path, common.DevParamsFilename, common.ParamsSchemaFilename, 1))
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
			obj[name] = FillDevParam(nestedProp, nil, nestedExampleValues)
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
