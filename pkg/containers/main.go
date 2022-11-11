package containers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/hasura/go-graphql-client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api"
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/rs/zerolog/log"
)

type ImageTag struct {
	Tag string `json:"tag"`
}

type BuildOptions struct {
	Tags []ImageTag
}

type RegistryManager struct {
	dockerClient ContainerClient
	graphClient  *graphql.Client
	// Package(string) error
}

type ContainerClient interface {
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImagePush(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error)
}

func NewRegistryManager() (*RegistryManager, error) {
	containerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	graphClient := api.NewClient()

	return &RegistryManager{
		dockerClient: containerClient,
		graphClient:  graphClient,
	}, nil
}

func (m *RegistryManager) Package(b *bundle.Bundle) error {
	imageURI := imageURIFromBundle(b)
	repository := repositoryFromBundle(b)
	log.Info().Msg(repository)
	opts := BuildOptions{
		Tags: []ImageTag{
			{
				Tag: imageURI,
			},
		},
	}

	errBuild := m.BuildImage(opts)
	if errBuild != nil {
		return errBuild
	}

	repoExists, errCheck := CheckRepositoryExists(repository)
	if errCheck != nil {
		return errCheck
	}

	if repoExists {
		return m.PushImage(imageURI)
	}

	log.Info().Msg("Repository does not exist, creating")
	return CreateRepository()
}

func (opts *BuildOptions) GetTags() []string {
	var tags []string
	for _, tag := range opts.Tags {
		tags = append(tags, tag.Tag)
	}
	return tags
}

// TODO: un-hack the static registry string
func imageURIFromBundle(b *bundle.Bundle) string {
	repository := "us-west1-docker.pkg.dev/md-wbeebe-0808-example-apps/sat-test-6789"
	return repository + "/" + b.Name + ":latest"
}

func repositoryFromBundle(b *bundle.Bundle) string {
	imageURI := imageURIFromBundle(b)
	return strings.Replace(imageURI, "/"+b.Name+":latest", "", -1)
}

func CheckRepositoryExists(repository string) (bool, error) {
	return true, nil
}

func CreateRepository() error {
	return nil
}

func (m *RegistryManager) BuildImage(opts BuildOptions) error {
	log.Info().Msg("Building image")
	ctx := context.TODO()

	// TODO: this is the "context" argument, make configurable
	tar, err := archive.TarWithOptions(".", &archive.TarOptions{})
	if err != nil {
		return err
	}

	cliOpts := types.ImageBuildOptions{
		// TODO: allow config from massdriver.yaml?
		Dockerfile: "./Dockerfile",
		Tags:       opts.GetTags(),
		// Remove: true,
	}

	res, errBuild := m.dockerClient.ImageBuild(ctx, tar, cliOpts)
	if errBuild != nil {
		return errBuild
	}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		logStr := strings.TrimSuffix(scanner.Text(), "\n")
		log := []byte(logStr)

		var logLine LogLine
		err := json.Unmarshal(log, &logLine)
		if err != nil {
			return err
		}
		msg := strings.TrimSuffix(logLine.Stream, "\n")
		if msg == "" {
			continue
		}
		fmt.Println(msg)
	}

	defer res.Body.Close()

	return nil
}

func (m *RegistryManager) ListImages() error {
	images, err := m.dockerClient.ImageList(context.TODO(), types.ImageListOptions{})
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE")

	for _, image := range images {
		if len(image.RepoDigests) == 0 || len(image.RepoTags) == 0 {
			continue
		}

		repo := strings.Split(image.RepoDigests[0], "@")[0]
		mostRecentTag := strings.Split(image.RepoTags[0], ":")[1]
		if repo == "<none>" {
			continue
		}
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", repo, mostRecentTag, image.ID[:10], "Created", humanFileSize(float64(image.VirtualSize))))
	}
	w.Flush()
	return nil
}

func (m *RegistryManager) PushImage(imageURI string) error {
	ctx := context.TODO()
	if m.graphClient == nil {
		return errors.New("graphClient is nil")
	}
	orgID := "ord-1"
	name := "sat-test-6789"

	authToken, err := api.GetToken(m.graphClient, orgID, name)
	if err != nil {
		return err
	}
	authStr, errCfg := getAuthConfig(authToken.Token, imageURI)
	if errCfg != nil {
		return errCfg
	}

	reader, errPush := m.dockerClient.ImagePush(ctx, imageURI, types.ImagePushOptions{
		// All           bool
		RegistryAuth: authStr, // RegistryAuth is the base64 encoded credentials for the registry
		// PrivilegeFunc RequestPrivilegeFunc
		// Platform      string
	})
	if errPush != nil {
		return errPush
	}
	defer reader.Close()
	return print(reader)
}
