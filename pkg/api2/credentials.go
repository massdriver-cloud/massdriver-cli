// Manages credential-type artifacts
package api2

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

var credentialArtifactDefinitions = []ArtifactDefinition{
	{"massdriver/aws-iam-role"},
	{"massdriver/azure-service-principal"},
	{"massdriver/gcp-service-account"},
	{"massdriver/kubernetes-cluster"},
}

// List supported credential types
func ListCredentialTypes() []ArtifactDefinition {
	return credentialArtifactDefinitions
}

// Get the first page of credentials for an artifac type
func ListCredentials(client graphql.Client, orgID string, artifacType string) ([]getArtifactsByTypeArtifactsPaginatedArtifactsItemsArtifact, error) {
	response, err := getArtifactsByType(context.Background(), client, orgID, artifacType)

	return response.Artifacts.Items, err
}

// Convert the API response to an Artifact
func (a *getArtifactsByTypeArtifactsPaginatedArtifactsItemsArtifact) ToArtifact() Artifact {
	return Artifact{
		ID:   a.Id,
		Name: a.Name,
	}
}
