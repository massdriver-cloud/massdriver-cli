package api2

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

func GetContainerRepository(client graphql.Client, artifactID, orgID, imageName, location string) (*ContainerRepository, error) {
	result := &ContainerRepository{}
	response, err := containerRepository(context.Background(), client, orgID, artifactID, ContainerRepositoryInput{ImageName: imageName, Location: location})
	if err != nil {
		return result, err
	}

	result.RepositoryUri = response.ContainerRepository.RepoUri
	result.Token = response.ContainerRepository.Token

	return result, nil
}
