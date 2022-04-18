package application

import (
	"encoding/json"
	"os"
	"path"
)

func GenerateFiles(baseDir string) error {
	applicationVariables := map[string]interface{}{
		"variable": map[string]interface{}{
			"connections": map[string]string{
				"type": "any",
			},
		},
	}

	err := os.MkdirAll(path.Join(baseDir, "src"), 0755)
	if err != nil {
		return err
	}

	applicationVariablesFile, err := os.Create(path.Join(baseDir, "src", "_application_variables.tf.json"))
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(applicationVariables, "", "  ")
	if err != nil {
		return err
	}
	_, err = applicationVariablesFile.Write(bytes)

	return err
}

// // Compile a JSON Schema to Terraform Variable Definition JSON
// func Compile(path string, out io.Writer) error {
// 	vars, err := getVars(path)
// 	if err != nil {
// 		return err
// 	}

// 	// You can't have an empty variable block, so if there are no vars return an empty json block
// 	if len(vars) == 0 {
// 		out.Write([]byte("{}"))
// 		return nil
// 	}

// 	variableFile := TFVariableFile{Variable: vars}

// 	bytes, err := json.MarshalIndent(variableFile, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	_, err = out.Write(bytes)

// 	return err
// }

// func getVars(path string) (map[string]TFVariable, error) {
// 	variables := map[string]TFVariable{}
// 	schema, err := jsonschema.GetJSONSchema(path)
// 	if err != nil {
// 		return variables, err
// 	}

// 	required := schema.Required

// 	for name, prop := range schema.Properties {
// 		variables[name] = NewTFVariable(prop, isRequired(name, required))
// 	}
// 	return variables, nil
// }

// func isRequired(name string, required []string) bool {
// 	for _, elem := range required {
// 		if name == elem {
// 			return true
// 		}
// 	}
// 	return false
// }
