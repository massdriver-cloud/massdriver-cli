package cdk8s

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
)

func GenerateFiles(bundlePath string, srcDir string) error {
	massdriverVariables := map[string]interface{}{
		"variable": map[string]interface{}{
			"md_metadata": map[string]string{
				"type": "any",
			},
		},
	}

	paramsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_params_variables.json"))
	if err != nil {
		return err
	}
	err = Compile(path.Join(bundlePath, common.ParamsSchemaFilename), paramsVariablesFile)
	if err != nil {
		return err
	}

	connectionsVariablesFile, err := os.Create(path.Join(bundlePath, srcDir, "_connections_variables.json"))
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
	devParamPath := path.Join(bundlePath, "src", common.DevParamsFilename)
	devParamsVariablesFile, err := os.OpenFile(devParamPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil { // fall back to create missing file
		devParamsVariablesFile, err = os.Create(devParamPath)
		if err != nil {
			return err
		}
	}

	err = CompileDevParams(devParamPath, devParamsVariablesFile)
	if err != nil {
		return fmt.Errorf("error compiling dev params: %w", err)
	}

	return nil
}
