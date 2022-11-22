package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/rs/zerolog/log"
)

const urlTemplate = "https://app.massdriver.cloud/projects/%s/targets/%v"

func DoDeployPreviewEnvironment(client graphql.Client, orgID string, id string, credentials []api2.Credential, previewConfig map[string]interface{}, ciContext map[string]interface{}) (*api2.Environment, error) {
	// TODO: target default connection binding by artifact
	// TODO: mutation signature
	log.Info().Str("project", id).Msg("Deploying preview environment.")

	templateData, err := json.Marshal(previewConfig)
	if err != nil {
		return nil, err
	}

	// buf := new(strings.Builder)
	// _, err = io.Copy(buf, templateData)

	// if err != nil {
	// 	return nil, err
	// }

	envVars := getOsEnv()
	// template := buf.String()
	config := os.Expand(string(templateData), func(s string) string { return envVars[s] })

	confMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(config), &confMap)

	previewEnv, err := api2.DeployPreviewEnvironment(client, orgID, id, credentials, confMap, ciContext)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlTemplate, id, previewEnv.ID)
	log.Info().
		Str("project", id).
		Str("url", url).
		Interface("environment", previewEnv.ID).
		Msg("Preview environment deploying.")
	exec.Command("open", url).Run()

	return &previewEnv, nil
}

func DeployPreviewEnvironment(client graphql.Client, orgID string, id string, previewConfigPath string, ciContextPath string) (*api2.Environment, error) {
	// TODO: parse ciContext
	ciContext := map[string]interface{}{
		"pull_request": map[string]interface{}{
			"title":  "Testing preview envs",
			"number": 9000,
		},
	}

	// TODO: read creds from conf file
	credentials := []api2.Credential{}

	previewConfigFile, err := os.Open(previewConfigPath)

	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(previewConfigFile)

	if err != nil {
		return nil, err
	}

	previewConfig := map[string]interface{}{}
	json.Unmarshal(byteValue, &previewConfig)

	return DoDeployPreviewEnvironment(client, orgID, id, credentials, previewConfig, ciContext)
}

func getOsEnv() map[string]string {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}

	osEnv := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	return osEnv
}
