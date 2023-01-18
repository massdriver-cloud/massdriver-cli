package api2

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

func GetContainerRepository(client graphql.Client, artifactID, orgID, imageName, location string) (*ContainerRepository, error) {
	containerRepository := &ContainerRepository{}
	response, err := getContainerRepository(context.Background(), client, artifactID, orgID, ContainerRepositoryInput{ImageName: imageName, Location: location})

	if err != nil {
		return containerRepository, err
	}

	containerRepository.RepositoryUri = response.ContainerRepository.RepoUri
	containerRepository.Token = response.ContainerRepository.Token

	return containerRepository, nil
}
