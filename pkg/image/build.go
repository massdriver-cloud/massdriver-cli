package image

import (
	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/rs/zerolog/log"
)

func Build(client graphql.Client, input PushImageInput, imageClient Client) error {
	log.Info().Str("image-name", input.ImageName).Str("platform", input.TargetPlatform).Msg("Starting to build the container image.")

	containerRepository, err := api2.GetContainerRepository(client, input.ArtifactId, input.OrganizationId, input.ImageName, input.Location)

	if err != nil {
		return err
	}

	res, err := imageClient.BuildImage(input, containerRepository)

	if err != nil {
		return err
	}

	err = handleResponseBuffer(res.Body)

	if err != nil {
		return err
	}
	log.Info().Str("image-name", input.ImageName).Msg("Successfully built!")

	return nil
}
