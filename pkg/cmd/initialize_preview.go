package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/massdriver-cloud/massdriver-cli/pkg/views/credential_types_table"
)

func InitializePreview(client graphql.Client, orgId string) (map[string]string, error) {
	selectedArtifactTypes, err := credential_types_table.New(api2.ListCredentialTypes())

	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUG] Selected ArtDefs: %v", selectedArtifactTypes)

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
	artifactIdMap := map[string]api2.Artifact{}
	artifactId := ""
	firstPage, err := api2.ListCredentials(client, orgId, artifacType)

	if err != nil {
		return artifactId, err
	}

	for _, artifactRecord := range firstPage {
		artifactIdMap[artifactRecord.Id] = artifactRecord.ToArtifact()
	}

	options := keys(artifactIdMap)

	if len(options) == 0 {
		// User has none of this type of artifact ... should this be an error?
		return "", nil
	}

	prompt := &survey.Select{
		Message: "Which credential?",
		Options: options,
		Description: func(value string, index int) string {
			artifactId := options[index]
			return artifactIdMap[artifactId].Name
		},
	}

	err = survey.AskOne(prompt, &artifactId)

	return artifactId, err
}

func keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
