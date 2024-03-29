package image

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

type ImageClient interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImagePush(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error)
}

type Client struct {
	Cli ImageClient
}

func NewImageClient() (Client, error) {
	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)

	if err != nil {
		return Client{}, errors.New("docker Engine API is not installed. to install it go to https://docs.docker.com/get-docker/ and follow the instructions")
	}

	return Client{Cli: cli}, nil
}

func (c *Client) BuildImage(input PushImageInput, containerRepository *api2.ContainerRepository) (*types.ImageBuildResponse, error) {
	tar, err := packageBuildDirectory(input.DockerBuildContext)

	if err != nil {
		return nil, err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: input.Dockerfile,
		Tags:       []string{imageFqn(containerRepository.RepositoryUri, input.ImageName, input.Tag)},
		Remove:     true,
		Platform:   input.TargetPlatform,
	}

	ctx := context.Background()

	res, err := c.Cli.ImageBuild(ctx, tar, opts)
	return &res, err
}

func (c *Client) PushImage(input PushImageInput, containerRepository *api2.ContainerRepository) (io.ReadCloser, error) {
	ctx := context.Background()
	auth, err := createAuthForCloud(containerRepository, input)

	if err != nil {
		return nil, err
	}

	res, err := c.Cli.ImagePush(ctx, imageFqn(containerRepository.RepositoryUri, input.ImageName, input.Tag), types.ImagePushOptions{RegistryAuth: auth})

	return res, err
}

func packageBuildDirectory(buildContext string) (io.ReadCloser, error) {
	return archive.TarWithOptions(buildContext, &archive.TarOptions{})
}

func imageFqn(uri, imageName, tag string) string {
	return fmt.Sprintf("%s/%s:%s", repoPrefix(uri), imageName, tag)
}

func repoPrefix(uri string) string {
	return strings.ReplaceAll(uri, "https://", "")
}

func createAuthForCloud(containerRepository *api2.ContainerRepository, input PushImageInput) (string, error) {
	authConfig := &types.AuthConfig{}

	err := setAuthUserNameByCloud(containerRepository, authConfig)

	if err != nil {
		return "", err
	}

	err = maybeRemoveSuffix(containerRepository, authConfig)

	if err != nil {
		return "", err
	}

	authConfig.Password = containerRepository.Token

	authConfigBytes, err := json.Marshal(authConfig)

	if err != nil {
		return "", err
	}

	encodedAuth := base64.URLEncoding.EncodeToString(authConfigBytes)

	return encodedAuth, nil
}

func setAuthUserNameByCloud(containerRepository *api2.ContainerRepository, auth *types.AuthConfig) error {
	switch identifyCloudByRepositoryUri(containerRepository.RepositoryUri) {
	case AWS:
		auth.Username = "AWS"
	case AZURE:
		auth.Username = "00000000-0000-0000-0000-000000000000"
	case GCP:
		auth.Username = "oauth2accesstoken"
	default:
		return fmt.Errorf("container repositories are not supported for %s", containerRepository.RepositoryUri)
	}

	return nil
}

func maybeRemoveSuffix(containerRepository *api2.ContainerRepository, auth *types.AuthConfig) error {
	r, err := regexp.Compile("[a-zA-Z0-9-_]+.docker.pkg.dev")

	if err != nil {
		return err
	}

	switch identifyCloudByRepositoryUri(containerRepository.RepositoryUri) {
	case GCP:
		auth.ServerAddress = r.FindString(containerRepository.RepositoryUri)
		return nil
	case AWS:
		auth.ServerAddress = containerRepository.RepositoryUri
		return nil
	case AZURE:
		auth.ServerAddress = containerRepository.RepositoryUri
		return nil
	default:
		return fmt.Errorf("container repositories are not supported for %s", containerRepository.RepositoryUri)
	}
}
