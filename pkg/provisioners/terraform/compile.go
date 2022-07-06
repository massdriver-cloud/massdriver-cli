package terraform

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
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
	err = Compile(path.Join(bundlePath, bundle.ParamsSchemaFilename), paramsVariablesFile)
	if err != nil {
		return err
	}

	connectionsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_connections_variables.tf.json"))
	if err != nil {
		return err
	}
	err = Compile(path.Join(bundlePath, bundle.ConnectionsSchemaFilename), connectionsVariablesFile)
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
	vars, err := getVars(path)
	if err != nil {
		return err
	}

	// You can't have an empty variable block, so if there are no vars return an empty json block
	if len(vars) == 0 {
		out.Write([]byte("{}"))
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

func getVars(path string) (map[string]TFVariable, error) {
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

func isRequired(name string, required []string) bool {
	for _, elem := range required {
		if name == elem {
			return true
		}
	}
	return false
}
