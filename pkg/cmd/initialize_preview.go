package cmd

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/massdriver-cloud/massdriver-cli/pkg/views/artifacts_table"
	"github.com/massdriver-cloud/massdriver-cli/pkg/views/credential_types_table"
)

// Present the initialize preview workflow
func InitializePreview(client graphql.Client, orgId string) (map[string]string, error) {
	selectedArtifactTypes, err := credential_types_table.New(api2.ListCredentialTypes())

	if err != nil {
		return nil, err
	}

	selectedCredentials := map[string]string{}

	for _, t := range selectedArtifactTypes {
		artifactId, err := initializePreviewPromptForCredentials(client, orgId, t.Name)
		if err != nil {
			return nil, err
		}
		selectedCredentials[t.Name] = artifactId
	}

	return selectedCredentials, nil
}

func initializePreviewPromptForCredentials(client graphql.Client, orgId string, artifacType string) (string, error) {
	artifactId := ""
	firstPage, err := api2.ListCredentials(client, orgId, artifacType)

	if err != nil {
		return artifactId, err
	}

	if len(firstPage) == 0 {
		// User has none of this type of artifact ... should this be an error?
		fmt.Printf("[INFO] No artifacts of type '%s' found.", artifacType)
		return "", nil
	}

	artifactList := []api2.Artifact{}
	for _, artifactRecord := range firstPage {
		// TODO: call ToArtifact() from ListCredentials
		artifactList = append(artifactList, artifactRecord.ToArtifact())
	}

	// TODO: set the table to only allowing one selection
	selectedArtifact, err := artifacts_table.New(artifactList)
	if err != nil {
		return artifactId, err
	}

	if len(selectedArtifact) == 0 {
		return artifactId, err
	}

	artifactId = selectedArtifact[0].ID

	return artifactId, nil
}
