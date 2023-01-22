package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
)

const urlTemplate = "https://app.massdriver.cloud/projects/%s/targets/%v"

type previewConfig struct {
	Credentials   map[string]string      `json:"credentials"`
	PackageParams map[string]interface{} `json:"packageParams"`
}

func (p *previewConfig) GetCredentials() []api2.Credential {
	credentials := []api2.Credential{}
	for k, v := range p.Credentials {
		cred := api2.Credential{
			ArtifactDefinitionType: k,
			ArtifactId:             v,
		}
		credentials = append(credentials, cred)
	}
	return credentials
}

func DoDeployPreviewEnvironment(client graphql.Client, orgID string, id string, credentials []api2.Credential, packageParams map[string]interface{}, ciContext map[string]interface{}) (*api2.Environment, error) {
	log.Info().Str("project", id).Msg("Deploying preview environment.")

	// interpolate template data
	templateData, err := json.Marshal(packageParams)
	if err != nil {
		return nil, err
	}

	envVars := getOsEnv()
	config := os.Expand(string(templateData), func(s string) string { return envVars[s] })

	interpolatedPackageParams := map[string]interface{}{}
	_ = json.Unmarshal([]byte(config), &interpolatedPackageParams)

	previewEnv, err := api2.DeployPreviewEnvironment(client, orgID, id, credentials, interpolatedPackageParams, ciContext)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlTemplate, previewEnv.Project.ID, previewEnv.ID)
	log.Info().
		Str("project", id).
		Str("url", url).
		Interface("environment", previewEnv.Slug).
		Msg("Preview environment deploying.")

	err = browser.OpenURL(url)
	_ = err
	return &previewEnv, nil
}

func DeployPreviewEnvironment(client graphql.Client, orgID string, id string, previewConfigPath string, ciContextPath string) (*api2.Environment, error) {
	ciContext := map[string]interface{}{}
	err := readJSONFile(ciContextPath, &ciContext)
	if err != nil {
		return nil, err
	}

	config := previewConfig{}
	err = readJSONFile(previewConfigPath, &config)

	if err != nil {
		return nil, err
	}

	return DoDeployPreviewEnvironment(client, orgID, id, config.GetCredentials(), config.PackageParams, ciContext)
}

func readJSONFile(filename string, v interface{}) error {
	fileBytes, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	err = json.Unmarshal(fileBytes, v)

	if err != nil {
		return err
	}

	return nil
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
