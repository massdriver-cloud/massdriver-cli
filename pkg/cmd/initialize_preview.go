package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/massdriver-cloud/massdriver-cli/pkg/config"
	"github.com/massdriver-cloud/massdriver-cli/pkg/views/artifacts_table"
	"github.com/massdriver-cloud/massdriver-cli/pkg/views/credential_types_table"
	"github.com/rs/zerolog/log"
)

// Present the initialize preview workflow
func InitializePreview(config *config.Config, projectSlugOrID string, previewCfgPath string) error {
	client := api2.NewClient(config.APIKey)
	previewCfg, err := DoInitializePreview(client, config.OrgID, projectSlugOrID)

	if err != nil {
		return err
	}

	return initializePreviewSerializeCfg(previewCfg, previewCfgPath)
}

const previewConfMode = 0600

func initializePreviewSerializeCfg(cfg map[string]interface{}, path string) error {
	previewConf, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, previewConf, previewConfMode)
	return err
}

func DoInitializePreview(client graphql.Client, orgID string, projectSlugOrID string) (map[string]interface{}, error) {
	defaultParams, err := initializePreviewGetProjectDefaultParams(client, orgID, projectSlugOrID)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get project")
		return nil, err
	}

	selectedArtifactTypes, err := credential_types_table.New(api2.ListCredentialTypes())

	if err != nil {
		log.Error().Err(err).Msg("Failed to get artifacts")
		return nil, err
	}

	selectedCredentials := map[string]string{}

	for _, t := range selectedArtifactTypes {
		artifactId, err := initializePreviewPromptForCredentials(client, orgID, t.Name)
		if err != nil {
			return nil, err
		}
		selectedCredentials[t.Name] = artifactId
	}

	conf := map[string]interface{}{
		"credentials":   selectedCredentials,
		"packageParams": defaultParams,
	}

	return conf, nil
}

func initializePreviewGetProjectDefaultParams(client graphql.Client, orgID string, projectSlugOrID string) (map[string]interface{}, error) {
	project, err := api2.GetProject(client, orgID, projectSlugOrID)
	if err != nil {
		return nil, err
	}

	return project.DefaultParams, nil
}

func initializePreviewPromptForCredentials(client graphql.Client, orgID string, artifacType string) (string, error) {
	artifactID := ""
	artifactList, err := api2.ListCredentials(client, orgID, artifacType)

	if err != nil {
		return artifactID, err
	}

	if len(artifactList) == 0 {
		// User has none of this type of artifact ... should this be an error?
		fmt.Printf("[INFO] No artifacts of type '%s' found.", artifacType)
		return "", nil
	}

	// TODO: set the table to only allowing one selection
	selectedArtifact, err := artifacts_table.New(artifactList)
	if err != nil {
		return artifactID, err
	}

	if len(selectedArtifact) == 0 {
		return artifactID, err
	}

	artifactID = selectedArtifact[0].ID

	return artifactID, nil
}
