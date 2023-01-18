package image

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/rs/zerolog/log"
)

const AWS = "AWS"
const GCP = "GCP"
const AZURE = "Azure"

func Push(client graphql.Client, input PushImageInput, imageClient Client) error {
	log.Info().Str("image-name", input.ImageName).Str("location", input.Location).Msg("Creating repository and fetching single use credentials")

	containerRepository, err := api2.GetContainerRepository(client, input.ArtifactId, input.OrganizationId, input.ImageName, input.Location)

	if err != nil {
		return err
	}

	cloudName := identifyCloudByRepositoryUri(containerRepository.RepositoryUri)

	log.Info().Str("cloud", cloudName).Msg("Credentials fetched successfully")
	log.Info().Str("tag", input.Tag).Msg("Building and tagging image")

	res, err := imageClient.BuildImage(input, containerRepository)

	if err != nil {
		return err
	}

	err = handleResponseBuffer(res.Body)

	if err != nil {
		return err
	}

	log.Info().Str("tag", input.Tag).Msg("Pushing image to repository. This may take a few minutes.")

	rd, err := imageClient.PushImage(input, containerRepository)

	if err != nil {
		return err
	}

	err = handleResponseBuffer(rd)

	if err != nil {
		return err
	}

	log.Info().Str("image", imageFqn(containerRepository.RepositoryUri, input.ImageName, input.Tag)).Msg("Image pushed successfully")

	return nil
}

func identifyCloudByRepositoryUri(uri string) string {
	switch {
	case strings.Contains(uri, "amazonaws.com"):
		return AWS
	case strings.Contains(uri, "azurecr.io"):
		return AZURE
	case strings.Contains(uri, "docker.pkg.dev"):
		return GCP
	default:
		return "unknown"
	}
}

func handleResponseBuffer(buf io.ReadCloser) error {
	defer buf.Close()

	return print(buf)
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
